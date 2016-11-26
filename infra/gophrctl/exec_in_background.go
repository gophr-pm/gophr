package main

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
)

func execInBackground(cmd string, args ...string) error {
	subProcess := exec.Command(cmd, args...)

	stdin, stdinErr := subProcess.StdinPipe()
	defer stdin.Close()

	subProcess.Stdout = os.Stdout
	subProcess.Stderr = os.Stderr

	startErr := subProcess.Start()
	waitErr := subProcess.Wait()

	if stdinErr != nil || startErr != nil || waitErr != nil {
		var b bytes.Buffer
		b.WriteString("Failed to exec \"")
		b.WriteString(cmd)
		for _, arg := range args {
			b.WriteByte(' ')
			b.WriteString(arg)
		}
		b.WriteString("\" due to some errors along the way: [ ")
		if stdinErr != nil {
			b.WriteString(stdinErr.Error())
			b.WriteString(", ")
		}
		if startErr != nil {
			b.WriteString(startErr.Error())
			b.WriteString(", ")
		}
		if waitErr != nil {
			b.WriteString(waitErr.Error())
		}
		b.WriteString(" ].")

		return errors.New(b.String())
	}

	return nil
}
