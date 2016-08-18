package main

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/briandowns/spinner"
	"gopkg.in/urfave/cli.v1"
)

var (
	loadingSpinner = spinner.New(spinner.CharSets[11], 100*time.Millisecond)
)

func readModule(command string, c *cli.Context) (module, error) {
	moduleID := c.Args().First()
	if len(moduleID) == 0 {
		// A blank module id means "all"
		moduleID = allModuleID
	}

	referencedModule := modules[moduleID]

	if referencedModule == nil {
		return nil, fmt.Errorf("\"%s\" is not a valid module", moduleID)
	}

	return referencedModule, nil
}

func newFailedModulesError(action string, failedModuleIds []string) error {
	// Create a buffer to list the modules that failed in the resulting error.
	buffer := bytes.Buffer{}
	buffer.WriteString("The following modules failed to ")
	buffer.WriteString(action)
	buffer.WriteByte(':')

	// Write out the list like [a,b,c,d] => "a, b, c and d".
	for i, failedModuleID := range failedModuleIds {
		if i > 0 {
			buffer.WriteString(", ")

			if i == (len(failedModuleIds) - 1) {
				buffer.WriteString("and ")
			}
		}

		buffer.WriteString(failedModuleID)
	}

	// Puncutate this message because we're civilized.
	buffer.WriteByte('.')

	return errors.New(buffer.String())
}

// TraverseModulesDependencyTree applies the iterator to a module and all of its
// dependencies in a depth-frst manner. NB: does not detect cycles.
func traverseModulesDependencyTree(
	moduleID string,
	visitedModules map[string]bool,
	iterator func(module),
) error {
	// Skip the "all" module since it isn't a real module.
	if moduleID == allModuleID {
		return nil
	}

	if !visitedModules[moduleID] {
		// Get the module that matches this module id.
		module := modules[moduleID]
		if module == nil {
			return fmt.Errorf("No such module \"%s\".", moduleID)
		}

		// Visit each dependency.
		for _, dependencyID := range module.deps() {
			err := traverseModulesDependencyTree(dependencyID, visitedModules, iterator)
			if err != nil {
				return err
			}
		}

		// Visit this module.
		iterator(module)
		visitedModules[moduleID] = true
	}

	return nil
}

// TraverseModules applies the iterator to all of the module and all of its
// their respective dependencies in such a way that no dependency is traversed
// before its dependant. NB: does not detect cycles.
func traverseModules(iterator func(module)) error {
	visitedModules := make(map[string]bool)

	for id := range modules {
		// Traverse the dependency tree of every module in order.
		if err := traverseModulesDependencyTree(id, visitedModules, iterator); err != nil {
			return err
		}
	}

	return nil
}

func modulesToString(excludeAll bool) string {
	var (
		buffer        bytes.Buffer
		isFirstModule = true
	)

	for moduleID := range modules {
		// Skip if this is the "all" module and we're supposed to skip it.
		if excludeAll && moduleID == allModuleID {
			continue
		}

		if !isFirstModule {
			buffer.WriteString(", ")
		} else {
			isFirstModule = false
		}

		buffer.WriteString(moduleID)
	}

	return buffer.String()
}

func exit(
	code int,
	c *cli.Context,
	command string,
	args ...interface{},
) {
	printError(args...)

	if c != nil {
		fmt.Println()
		if len(command) > 0 {
			cli.ShowCommandHelp(c, command)
		} else {
			cli.ShowAppHelp(c)
		}
	}

	os.Exit(code)
}

func startSpinner(message string) {
	loadingSpinner.Color("green")
	loadingSpinner.Suffix = " " + message
	loadingSpinner.FinalMSG = "done."
	loadingSpinner.Start()
}

func stopSpinner() {
	loadingSpinner.Stop()
	fmt.Println()
}
