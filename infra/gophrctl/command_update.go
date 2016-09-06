package main

import (
	"fmt"
	"path/filepath"

	"gopkg.in/urfave/cli.v1"
)

func updateCommand(c *cli.Context) error {
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
		printInfo("Updating all modules")
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

			if err = updateModule(m, gophrRoot, env); err != nil {
				goto exitWithError
			}
		}
		printSuccess("All modules were updated successfully")
	} else if m, exists = modules[moduleName]; exists {
		printInfo(fmt.Sprintf("Updating module \"%s\"", moduleName))
		if env == environmentDev {
			if err = assertMinikubeRunning(); err != nil {
				goto exitWithError
			}
		}

		if err = updateModule(m, gophrRoot, env); err != nil {
			goto exitWithError
		}
		printSuccess(fmt.Sprintf("Module \"%s\" was updated successfully", moduleName))
	} else {
		err = newNoSuchModuleError(moduleName)
		goto exitWithErrorAndHelp
	}

	return nil
exitWithError:
	exit(exitCodeUpdateFailed, nil, "", err)
	return nil
exitWithErrorAndHelp:
	exit(exitCodeUpdateFailed, c, "update", err)
	return nil
}

func updateModule(m *module, gophrRoot string, env environment) error {
	// Apply in order.
	for _, k8sfile := range m.k8sfiles {
		// Put together the absolute path.
		k8sfilePath := filepath.Join(gophrRoot, fmt.Sprintf("%s.%s.yml", k8sfile, env))
		// Perform the create command.
		if err := applyInK8S(k8sfilePath); err != nil {
			return err
		}
	}

	return nil
}
