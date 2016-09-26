package main

import (
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
		if excludedStr := c.String(flagNameExclude); len(excluded) > 0 {
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
		var prevModuleName string
		for _, m := range orderedModules {
			// Wait for the previous module before starting this one.
			if len(prevModuleName) > 0 {
				if err = waitForK8SPods(prevModuleName); err != nil {
					return err
				}
			}

			// Use the environment to toggle the unfiltered list.
			var k8sfiles []string
			if env == environmentProd {
				k8sfiles = m.prodK8SFiles
			} else {
				k8sfiles = m.devK8SFiles
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

			// Set the prev module name so on the next iteration we can wait for it.
			prevModuleName = m.name
		}

		// We only get here if everything worked out.
		printSuccess("Modules brought up successfully")
		return nil
	}); err != nil {
		exit(exitCodeUpFailed, nil, "", err)
	}

	return nil
}
