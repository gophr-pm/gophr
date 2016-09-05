package main

import (
	"fmt"
	"path/filepath"

	"gopkg.in/urfave/cli.v1"
)

func cycleCommand(c *cli.Context) error {
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
		printInfo("Cycling all modules")
		if err = assertMinikubeRunning(); err != nil {
			goto exitWithError
		}
		for _, m = range modules {
			// Check for db inclusion.
			if m.name == "db" && !c.Bool(flagNameIncludeDB) {
				continue
			}

			if err = cycleModule(m, gophrRoot, env); err != nil {
				goto exitWithError
			}
		}
		printSuccess("All modules were cycled successfully")
	} else if m, exists = modules[moduleName]; exists {
		printInfo(fmt.Sprintf("Cycling module \"%s\"", moduleName))
		if err = assertMinikubeRunning(); err != nil {
			goto exitWithError
		}
		if err = cycleModule(m, gophrRoot, env); err != nil {
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

func cycleModule(m *module, gophrRoot string, env environment) error {
	// Destroy in reverse order.
	for i := len(m.k8sfiles) - 1; i >= 0; i-- {
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
