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
	var (
		m          *module
		err        error
		env        environment
		exists     bool
		gophrRoot  string
		moduleName string
	)

	if env, err = readEnvironment(c); err != nil {
		goto exitWithError
	}

	if gophrRoot, err = readGophrRoot(c); err != nil {
		goto exitWithError
	}

	moduleName = c.Args().First()
	if len(moduleName) == 0 {
		// Means "all modules".
		printInfo("Building all modules")
		if err = assertMinikubeRunning(); err != nil {
			goto exitWithError
		}
		for _, m = range modules {
			if err = buildModule(m, gophrRoot, env); err != nil {
				goto exitWithError
			}
		}
		printSuccess("All modules were built successfully")
	} else if m, exists = modules[moduleName]; exists {
		printInfo(fmt.Sprintf("Building module \"%s\"", moduleName))
		if err = assertMinikubeRunning(); err != nil {
			goto exitWithError
		}
		if err = buildModule(m, gophrRoot, env); err != nil {
			goto exitWithError
		}
		printSuccess(fmt.Sprintf("Module \"%s\" was built successfully", moduleName))
	} else {
		err = newNoSuchModuleError(moduleName)
		goto exitWithErrorAndHelp
	}

	return nil
exitWithError:
	exit(exitCodeBuildFailed, nil, "", err)
	return nil
exitWithErrorAndHelp:
	exit(exitCodeBuildFailed, c, "build", err)
	return nil
}

func buildModule(m *module, gophrRoot string, env environment) error {
	buildArgs := buildInMinikubeArgs{
		imageTag:       devDockerImageTag, // TODO(skeswa): tag should depend on env.
		imageName:      fmt.Sprintf("gophr-%s-%s", m.name, env),
		contextPath:    filepath.Join(gophrRoot, m.buildContext),
		dockerfilePath: filepath.Join(gophrRoot, fmt.Sprintf("%s.%s", m.dockerfile, env)),
	}

	return buildInMinikube(buildArgs)
}
