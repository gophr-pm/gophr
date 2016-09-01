package main

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/urfave/cli.v1"
)

const (
	cliVersion = "0.0.1"
	envTypeDev = "dev"

	flagNameEnv        = "env"
	flagNameRepoPath   = "repo-path"
	flagNameForeground = "foreground"

	commandNameBuild   = "build"
	commandDescBuild   = "Updates module images"
	commandNameLog     = "log"
	commandDescLog     = "Logs module's container output to stdout"
	commandNameRestart = "restart"
	commandDescRestart = "Restarts module containers"
	commandNameSSH     = "ssh"
	commandDescSSH     = "Starts a shell session within a module's container"
	commandNameStart   = "start"
	commandDescStart   = "Start module containers"
	commandNameStop    = "stop"
	commandDescStop    = "Stops module containers"
	commandNameTest    = "test"
	commandDescTest    = "Runs module tests"
)

var (
	moduleCommandArgsUsage           = fmt.Sprintf("[%s] [arguments...]", modulesToString())
	moduleCommandArgsUsageWithoutAll = fmt.Sprintf("[%s] [arguments...]", modulesToString())
)

func main() {
	var (
		app    = cli.NewApp()
		gopath = os.Getenv("GOPATH")

		defaultRepoPath string
	)

	// First, make sure that there is a go path.
	if len(gopath) < 1 {
		exit(1, nil, "", "$GOPATH must be defined for gophrctl to work.")
	}

	// Next, create the default repo path
	defaultRepoPath = filepath.Join(gopath, "src/github.com/skeswa/gophr")

	// Then, describe command metadata.
	app.Name = "gophrctl"
	app.Usage = "Manages the gophr development and deployment environment."
	app.Version = cliVersion
	app.HelpName = "gophrctl"

	// After that, set the global flags for gophrctl.
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  flagNameEnv,
			Value: envTypeDev,
			Usage: "gophr execution environment",
		},
		cli.StringFlag{
			Name:  flagNameRepoPath,
			Value: defaultRepoPath,
			Usage: "path to the gophr repository",
		},
	}

	// Next, reference every command capable to gophrctl.
	app.Commands = []cli.Command{
		// Build command.
		{
			Name:      commandNameBuild,
			Usage:     commandDescBuild,
			Action:    buildCommand,
			ArgsUsage: moduleCommandArgsUsage,
		},

		// Log command.
		{
			Name:      commandNameLog,
			Usage:     commandDescLog,
			Action:    logCommand,
			ArgsUsage: moduleCommandArgsUsageWithoutAll,
		},

		// Restart command.
		{
			Name:      commandNameRestart,
			Usage:     commandDescRestart,
			Action:    restartCommand,
			ArgsUsage: moduleCommandArgsUsage,
		},

		// SSH command.
		{
			Name:      commandNameSSH,
			Usage:     commandDescSSH,
			Action:    sshCommand,
			ArgsUsage: moduleCommandArgsUsageWithoutAll,
		},

		// Start command.
		{
			Name:      commandNameStart,
			Usage:     commandDescStart,
			Action:    startCommand,
			ArgsUsage: moduleCommandArgsUsage,

			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  flagNameForeground + ",f",
					Usage: "makes the container run in the foreground",
				},
			},
		},

		// Stop command.
		{
			Name:      commandNameStop,
			Usage:     commandDescStop,
			Action:    stopCommand,
			ArgsUsage: moduleCommandArgsUsage,
		},

		// Test command.
		{
			Name:      commandNameTest,
			Usage:     commandDescTest,
			Action:    testCommand,
			ArgsUsage: moduleCommandArgsUsage,
		},
	}

	// Lastly, execute the command line application.
	app.Run(os.Args)
}
