package main

import (
	"fmt"
	"path/filepath"
	"strings"

	"gopkg.in/urfave/cli.v1"
)

func cycleCommand(c *cli.Context) error {
	var (
		m          *module
		err        error
		env        = readEnvironment(c)
		exists     bool
		gophrRoot  string
		moduleName string
	)

	if gophrRoot, err = readGophrRoot(c); err != nil {
		goto exitWithError
	}

	moduleName = c.Args().First()
	if len(moduleName) == 0 {
		// Means "all modules".
		printInfo("Cycling all modules")
		if env == environmentDev {
			if err = assertMinikubeRunning(); err != nil {
				goto exitWithError
			}
		}

		for _, m = range modules {
			// Check for db inclusion.
			if m.name == "db" && !c.Bool(flagNameIncludeDB) {
				continue
			}

			if err = cycleModule(c, m, gophrRoot, env); err != nil {
				goto exitWithError
			}
		}
		printSuccess("All modules were cycled successfully")
	} else if m, exists = modules[moduleName]; exists {
		printInfo(fmt.Sprintf("Cycling module \"%s\"", moduleName))
		if env == environmentDev {
			if err = assertMinikubeRunning(); err != nil {
				goto exitWithError
			}
		}

		if err = cycleModule(c, m, gophrRoot, env); err != nil {
			goto exitWithError
		}
		printSuccess(fmt.Sprintf("Module \"%s\" was cycled successfully", moduleName))
	} else {
		err = newNoSuchModuleError(moduleName)
		goto exitWithErrorAndHelp
	}

	return nil
exitWithError:
	exit(exitCodeCycleFailed, nil, "", err)
	return nil
exitWithErrorAndHelp:
	exit(exitCodeCycleFailed, c, "cycle", err)
	return nil
}

func cycleModule(c *cli.Context, m *module, gophrRoot string, env environment) error {
	if env == environmentProd {
		var (
			err               error
			k8sProdContext    string
			oldK8SProdContext string
		)

		// Read the production context before continuing.
		if k8sProdContext, err = readK8SProdContext(c); err != nil {
			return err
		}

		// Switch to the production context then switch back afterwards.
		if oldK8SProdContext, err = switchK8SContext(k8sProdContext); err != nil {
			return err
		}
		defer switchK8SContext(oldK8SProdContext)
	}

	// Memorize whether services should be deleted.
	shouldDeleteServices := c.Bool(flagNameDeleteServices)

	// Destroy in reverse order.
	for i := len(m.k8sfiles) - 1; i >= 0; i-- {
		// Only delete services if that flag says so.
		if strings.HasSuffix(m.k8sfiles[i], "service") && !shouldDeleteServices {
			continue
		}

		// Put together the absolute path.
		k8sfile := m.k8sfiles[i]
		k8sfilePath := filepath.Join(gophrRoot, fmt.Sprintf("%s.%s.yml", k8sfile, env))
		// Only destroy if its already a thing.
		if existsInK8S(k8sfilePath) {
			if err := deleteInK8S(k8sfilePath); err != nil {
				return err
			}
		}
	}

	// Create in order.
	for _, k8sfile := range m.k8sfiles {
		// Put together the absolute path.
		k8sfilePath := filepath.Join(gophrRoot, fmt.Sprintf("%s.%s.yml", k8sfile, env))
		// Perform the create command.
		if err := createInK8S(k8sfilePath); err != nil {
			return err
		}
	}

	return nil
}
