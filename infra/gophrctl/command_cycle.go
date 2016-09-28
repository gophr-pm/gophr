package main

import (
	"fmt"

	"gopkg.in/urfave/cli.v1"
)

func cycleCommand(c *cli.Context) error {
	if err := runInK8S(c, func() error {
		var (
			env = readEnvironment(c)

			m          *module
			err        error
			exists     bool
			gophrRoot  string
			moduleName string
		)

		// First, let's get ourselves oriented.
		if gophrRoot, err = readGophrRoot(c); err != nil {
			return err
		}

		moduleName = c.Args().First()
		if len(moduleName) == 0 {
			// Means "all modules".
			printInfo("Cycling all modules")
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

				if err = cycleModule(c, m, gophrRoot, env); err != nil {
					return err
				}
			}
			printSuccess("All modules were cycled successfully")
		} else if m, exists = modules[moduleName]; exists {
			printInfo(fmt.Sprintf("Cycling module \"%s\"", moduleName))
			if env == environmentDev {
				if err = assertMinikubeRunning(); err != nil {
					return err
				}
			}

			if err = cycleModule(c, m, gophrRoot, env); err != nil {
				return err
			}
			printSuccess(fmt.Sprintf("Module \"%s\" was cycled successfully", moduleName))
		} else {
			err = newNoSuchModuleError(moduleName)
			return err
		}

		return nil
	}); err != nil {
		exit(exitCodeCycleFailed, nil, "", err)
	}

	return nil
}

func cycleModule(c *cli.Context, m *module, gophrRoot string, env environment) error {
	// Memorize whether services should be deleted.
	shouldDeletePersistent := c.Bool(flagNameDeletePersistent)

	// Make sure that volumes exist first.
	if env == environmentProd && len(m.prodVolumes) > 0 {
		if err := createGCloudVolumesIfNotExist(m.prodVolumes...); err != nil {
			return err
		}
	}

	// Use the environment to toggle the unfiltered list.
	k8sfilePaths, err := getModuleK8SFilePaths(c, m)
	if err != nil {
		return err
	}
	// Make sure any potential generated files get deleted.
	defer deleteGeneratedK8SFiles(k8sfilePaths)

	// Destroy in reverse order.
	for i := len(k8sfilePaths) - 1; i >= 0; i-- {
		k8sfilePath := k8sfilePaths[i]

		// Only delete persistent resources if that flag says so.
		if !shouldDeletePersistent && isPersistentK8SResource(k8sfilePath) {
			continue
		}

		// Only destroy if its already a thing.
		if existsInK8S(k8sfilePath) {
			if err := deleteInK8S(k8sfilePath); err != nil {
				return err
			}
		}
	}

	// Create in order.
	for _, k8sfilePath := range k8sfilePaths {
		// Perform the create command.
		if !existsInK8S(k8sfilePath) {
			if err := createInK8S(k8sfilePath); err != nil {
				return err
			}
		}
	}

	return nil
}
