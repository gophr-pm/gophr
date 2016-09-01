package main

import "fmt"

type execError struct {
	err    error
	output []byte
}

func newExecError(output []byte, err error) *execError {
	return &execError{err: err, output: output}
}

func (e *execError) Error() string {
	return fmt.Sprintf("Command execution failed: %v.\n\n%s\n", e.err, string(e.output[:]))
}

func newNoSuchModuleError(nonExistentModuleName string) error {
	return fmt.Errorf("Could not find a module called \"%s\"", nonExistentModuleName)
}
