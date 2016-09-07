package main

import (
	"fmt"
	"path/filepath"
	"strings"

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
	shouldDeleteServices := c.Bool(flagNameDeleteServices)

	// Filter the k8sfiles.
	var k8sfiles []string
	for _, k8sfile := range m.k8sfiles {
		// Only delete services if that flag says so.
		if strings.HasSuffix(k8sfile, "service") && !shouldDeleteServices {
			continue
		}
		// Ignore storage.
		if strings.HasSuffix(k8sfile, "storage") {
			continue
		}

		k8sfiles = append(k8sfiles, k8sfile)
	}

	// Destroy in reverse order.
	for i := len(k8sfiles) - 1; i >= 0; i-- {
		// Put together the absolute path.
		k8sfilePath := filepath.Join(gophrRoot, fmt.Sprintf("%s.%s.yml", k8sfiles[i], env))
		// Only destroy if its already a thing.
		if existsInK8S(k8sfilePath) {
			if err := deleteInK8S(k8sfilePath); err != nil {
				return err
			}
		}
	}

	// Create in order.
	for _, k8sfile := range k8sfiles {
		// Put together the absolute path.
		k8sfilePath := filepath.Join(gophrRoot, fmt.Sprintf("%s.%s.yml", k8sfile, env))
		// Perform the create command.
		if err := createInK8S(k8sfilePath); err != nil {
			return err
		}
	}

	return nil
}
