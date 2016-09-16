package main

import (
	"fmt"
	"path/filepath"

	"gopkg.in/urfave/cli.v1"
)

func updateCommand(c *cli.Context) error {
	if err := runInK8S(c, func() error {
		var (
			m          *module
			err        error
			env        = readEnvironment(c)
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
			printInfo("Updating all modules")
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

				if err = updateModule(c, m, gophrRoot, env); err != nil {
					return err
				}
			}
			printSuccess("All modules were updated successfully")
		} else if m, exists = modules[moduleName]; exists {
			printInfo(fmt.Sprintf("Updating module \"%s\"", moduleName))
			if env == environmentDev {
				if err = assertMinikubeRunning(); err != nil {
					return err
				}
			}

			if err = updateModule(c, m, gophrRoot, env); err != nil {
				return err
			}
			printSuccess(fmt.Sprintf("Module \"%s\" was updated successfully", moduleName))
		} else {
			err = newNoSuchModuleError(moduleName)
			return err
		}

		return nil
	}); err != nil {
		exit(exitCodeUpdateFailed, nil, "", err)
	}

	return nil
}

func updateModule(c *cli.Context, m *module, gophrRoot string, env environment) error {
	// Use the environment to toggle the unfiltered list.
	var k8sfiles []string
	if env == environmentProd {
		k8sfiles = m.prodK8SFiles
	} else {
		k8sfiles = m.devK8SFiles
	}

	// Apply in order.
	for _, k8sfile := range k8sfiles {
		// Put together the absolute path.
		k8sfilePath := filepath.Join(gophrRoot, k8sfile)
		// Perform the create command.
		if err := applyInK8S(k8sfilePath); err != nil {
			return err
		}
	}

	return nil
}
