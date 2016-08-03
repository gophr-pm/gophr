package main

import (
	"bytes"
	"errors"
	"fmt"

	"gopkg.in/urfave/cli.v1"
)

func readModule(c *cli.Context) (module, error) {
	// moduleID := c.Args().First()
	// referencedModule := modules[moduleID]
	return nil, nil
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
	var visitedModules map[string]bool

	for id := range modules {
		if err := traverseModulesDependencyTree(id, visitedModules, iterator); err != nil {
			return err
		}
	}

	return nil
}
