package main

import (
	"fmt"
	"path/filepath"
	"strings"

	cli "gopkg.in/urfave/cli.v1"
)

func upCommand(c *cli.Context) error {
	if err := runInK8S(c, func() error {
		// Scope out the current environment.
		env := readEnvironment(c)
		// Next, let's get ourselves oriented.
		gophrRoot, err := readGophrRoot(c)
		if err != nil {
			return err
		}
		// Lastly, check out what should be excluded.
		var excluded map[string]bool
		if excludedStr := c.String(flagNameExclude); len(excludedStr) > 0 {
			excluded = make(map[string]bool)
			excludedStrParts := strings.Split(excludedStr, ",")
			for _, excludedStrPart := range excludedStrParts {
				excluded[excludedStrPart] = true
			}
		}

		// Order the modules so that things start with everything ready to go.
		printInfo("Bringing modules up")
		orderedModules := orderModulesByDeps(modules, excluded, false)

		// For each module, in order, wait for the previous module before starting.
		for _, m := range orderedModules {
			// Do docker build in dev if not exists.
			if env == environmentDev {
				var (
					imageTag  = devDockerImageTag
					imageName = fmt.Sprintf("gophr-%s-%s", m.name, env)
				)

				if built, buildCheckErr := isBuiltInMinikube(
					imageName,
					imageTag,
				); buildCheckErr != nil {
					return err
				} else if !built {
					if err = buildInMinikube(buildInMinikubeArgs{
						imageTag:    imageTag,
						imageName:   imageName,
						contextPath: filepath.Join(gophrRoot, m.buildContext),
						dockerfilePath: filepath.Join(
							gophrRoot,
							fmt.Sprintf("%s.%s", m.dockerfile, env)),
					}); err != nil {
						return err
					}
				}
			}

			// Use the environment to toggle the unfiltered list.
			k8sfiles, err := getModuleK8SFilePaths(c, m)
			if err != nil {
				return err
			}
			// Make sure any potential generated files get deleted.
			defer deleteGeneratedK8SFiles(k8sfiles)

			// Delete module in k8s only if transient (can exit).
			if m.transient {
				// Delete in order (if exists).
				for _, k8sfile := range k8sfiles {
					// Put together the absolute path.
					k8sfilePath := filepath.Join(gophrRoot, k8sfile)
					// Perform the delete command.
					if existsInK8S(k8sfilePath) {
						if err = deleteInK8S(k8sfilePath); err != nil {
							return err
						}
					}
				}
			}

			// Create in order (if not exists).
			for _, k8sfile := range k8sfiles {
				// Put together the absolute path.
				k8sfilePath := filepath.Join(gophrRoot, k8sfile)
				// Perform the create command.
				if !existsInK8S(k8sfilePath) {
					if err = createInK8S(k8sfilePath); err != nil {
						return err
					}
				}
			}

			// Wait for the module to start (or finish) before continuing.
			waitTilFinished := m.transient
			if err = waitForK8SPods(m.name, waitTilFinished); err != nil {
				return err
			}
		}

		// We only get here if everything worked out.
		printSuccess("Modules brought up successfully")
		return nil
	}); err != nil {
		exit(exitCodeUpFailed, nil, "", err)
	}

	return nil
}
