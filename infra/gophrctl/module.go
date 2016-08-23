package main

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"gopkg.in/urfave/cli.v1"
)

func getModules(c *cli.Context) ([]string, error) {
	imagesPath := filepath.Join(c.GlobalString(flagNameRepoPath), "infra/images")
	files, err := ioutil.ReadDir(imagesPath)
	if err != nil {
		return nil, err
	}

	var modules []string
	for _, file := range files {
		if file.IsDir() {
			modules = append(modules, file.Name())
		}
	}

	return modules, nil
}

func readModule(c *cli.Context, modules []string) (string, error) {
	proposedModule := c.Args().First()

	if isEveryModule(proposedModule) {
		return proposedModule, nil
	}

	for _, module := range modules {
		if module == proposedModule {
			return module, nil
		}
	}

	return "", fmt.Errorf("Invalid module \"%s\".", proposedModule)
}

func isEveryModule(module string) bool {
	return len(module) < 1
}

//
// type module interface {
// 	id() string
// 	build(*cli.Context) error
// 	start(*cli.Context) error
// 	stop(*cli.Context) error
// 	log(*cli.Context) error
// 	ssh(*cli.Context) error
// 	test(*cli.Context) error
// 	restart(*cli.Context) error
// }
//
// var modules = map[string]module{
// 	allModuleID:     &allModule{},
// 	apiModuleID:     &apiModule{baseModule{apiModuleID}},
// 	dbModuleID:      &dbModule{baseModule{dbModuleID}},
// 	indexerModuleID: &indexerModule{baseModule{indexerModuleID}},
// 	routerModuleID:  &routerModule{baseModule{routerModuleID}},
// 	webModuleID:     &webModule{baseModule{webModuleID}},
// }
//
// func doModuleBuild(
// 	moduleID string,
// 	targetDev bool,
// 	exitOnError bool,
// 	workDir string,
// ) error {
// 	printInfo("Building", moduleID+".")
//
// 	// Perform the docker build.
// 	startSpinner("Executing docker build...")
// 	err := doDockerBuild(
// 		workDir,
// 		filepath.Join(workDir, modules[moduleID].dockerfile()),
// 		dockerImageNameOf(moduleID),
// 		dockerDevImageTag,
// 	)
// 	stopSpinner()
//
// 	// Report on results.
// 	if err != nil {
// 		printError("Failed to build", moduleID+":")
// 		print(err)
//
// 		// Only exit if necessary.
// 		if exitOnError {
// 			os.Exit(exitCodeBuildFailed)
// 		}
// 	} else {
// 		printSuccess("Built", moduleID, "successfully.")
// 	}
//
// 	return nil
// }
//
// func doModuleStart(
// 	moduleID string,
// 	targetDev bool,
// 	exitOnError bool,
// 	workDir string,
// 	backgrounded bool,
// ) error {
// 	printInfo("Starting", moduleID+".")
//
// 	// Localize container metadata.
// 	ports, links, volumes := modules[moduleID].containerMetadata()
//
// 	// Perform the docker build.
// 	if backgrounded {
// 		startSpinner("Executing docker run...")
// 	}
// 	err := doDockerRun(
// 		workDir,
// 		dockerImageNameOf(moduleID),
// 		dockerDevImageTag,
// 		dockerContainerNameOf(moduleID),
// 		backgrounded,
// 		ports,
// 		links,
// 		volumes,
// 	)
// 	if backgrounded {
// 		stopSpinner()
// 	}
//
// 	// Report on results.
// 	if err != nil {
// 		printError("Failed to start", moduleID+":")
// 		print(err)
//
// 		// Only exit if necessary.
// 		if exitOnError {
// 			os.Exit(exitCodeStartFailed)
// 		}
// 	} else {
// 		printSuccess("Started", moduleID, "successfully.")
// 	}
//
// 	return nil
// }
