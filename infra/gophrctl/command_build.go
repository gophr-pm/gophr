package main

import (
	"fmt"
	"path/filepath"

	"gopkg.in/urfave/cli.v1"
)

const (
	devDockerImageTag = "v1"
)

func buildCommand(c *cli.Context) error {
	if len(c.Args().First()) == 0 {
		printInfo("Building all modules")
	}

	if err := runInK8S(c, func() error {
		var (
			env = readEnvironment(c)

			m          *module
			err        error
			exists     bool
			gophrRoot  string
			moduleName string
		)

		if gophrRoot, err = readGophrRoot(c); err != nil {
			return err
		}

		moduleName = c.Args().First()
		if len(moduleName) == 0 {
			// Means "all modules".
			if env == environmentDev {
				if err = assertMinikubeRunning(); err != nil {
					return err
				}
			}

			for _, m = range modules {
				// Check for db inclusion.
				if m.name == "db" && !c.Bool(flagNameIncludeDB) {
					continue
				}

				if err = buildModule(c, m, gophrRoot, env); err != nil {
					return err
				}
			}
			printSuccess("All modules were built successfully")
		} else if m, exists = modules[moduleName]; exists {
			printInfo(fmt.Sprintf("Building module \"%s\"", moduleName))
			if env == environmentDev {
				if err = assertMinikubeRunning(); err != nil {
					return err
				}
			}

			if err = buildModule(c, m, gophrRoot, env); err != nil {
				return err
			}
			printSuccess(fmt.Sprintf("Module \"%s\" was built successfully", moduleName))
		} else {
			err = newNoSuchModuleError(moduleName)
			return err
		}

		return nil
	}); err != nil {
		exit(exitCodeBuildFailed, nil, "", err)
	}

	return nil
}

func buildModule(c *cli.Context, m *module, gophrRoot string, env environment) error {
	switch env {
	case environmentDev:
		return buildInMinikube(buildInMinikubeArgs{
			imageTag:       devDockerImageTag,
			imageName:      fmt.Sprintf("gophr-%s-%s", m.name, env),
			contextPath:    filepath.Join(gophrRoot, m.buildContext),
			dockerfilePath: filepath.Join(gophrRoot, fmt.Sprintf("%s.%s", m.dockerfile, env)),
		})
	case environmentProd:
		var (
			err       error
			gpi       string
			version   imageVersion
			imageName = "gophr-" + m.name
		)

		// Get the google project id for the push.
		if gpi, err = readGPI(c); err != nil {
			return err
		}

		// Bump the version in the versionfile.
		if version, err = promptImageVersionBump(filepath.Join(gophrRoot, m.versionfile)); err != nil {
			return err
		}

		// Build the prod docker image outside of minikube, so as not to crowd it.
		if err = dockerBuild(dockerBuildArgs{
			gpi:            gpi,
			latest:         true,
			imageTag:       version.String(),
			imageName:      imageName,
			contextPath:    filepath.Join(gophrRoot, m.buildContext),
			dockerfilePath: filepath.Join(gophrRoot, fmt.Sprintf("%s.%s", m.dockerfile, env)),
		}); err != nil {
			return err
		}

		// Push the newly built image to gcr.
		if err = dockerPush(gpi, imageName, version.String()); err != nil {
			return err
		}

		// Update the kubernetes configuration files.
		startSpinner("Updating kubernetes configuration")
		for _, k8sfile := range m.k8sfiles {
			var (
				newImageURL = fmt.Sprintf("gcr.io/%s/%s:%s", gpi, imageName, version.String())
				k8sfilePath = filepath.Join(gophrRoot, fmt.Sprintf("%s.%s.yml", k8sfile, env))
			)

			if err = updateProdK8SFileImage(newImageURL, k8sfilePath); err != nil {
				stopSpinner(false)
				return err
			}
		}
		stopSpinner(true)
	}

	return nil
}
