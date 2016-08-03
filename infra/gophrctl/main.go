package main

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/urfave/cli.v1"
)

const (
	cliVersion = "0.0.1"

	flagNameEnv      = "env"
	flagNameRepoPath = "repo-path"

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
	moduleCommandArgsUsage           = fmt.Sprintf("[module (%s)] [arguments...]", modulesToString(false))
	moduleCommandArgsUsageWithoutAll = fmt.Sprintf("[module (%s)] [arguments...]", modulesToString(true))
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
	if strings.HasSuffix(gopath, "/") {
		defaultRepoPath = gopath + "src/github.com/skeswa/gophr"
	} else {
		defaultRepoPath = gopath + "/src/github.com/skeswa/gophr"
	}

	// Then, describe command metadata.
	app.Name = "gophrctl"
	app.Usage = "Manages the gophr development and deployment environment."
	app.Version = cliVersion
	app.HelpName = "gophrctl"

	// After that, set the global flags for gophrctl.
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  flagNameEnv,
			Value: "dev",
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
			ArgsUsage: moduleCommandArgsUsage,

			Action: func(c *cli.Context) error {
				fmt.Println(c.GlobalString(flagNameRepoPath))
				// Read the module.
				m, err := readModule(commandNameBuild, c)
				if err != nil {
					exit(1, c, commandNameBuild, err)
				}

				m.build(c, false)
				return nil
			},
		},

		// Log command.
		{
			Name:      commandNameLog,
			Usage:     commandDescLog,
			ArgsUsage: moduleCommandArgsUsageWithoutAll,

			Action: func(c *cli.Context) error {
				// Read the module.
				m, err := readModule(commandNameLog, c)
				if err != nil {
					return err
				}

				return m.build(c, false)
			},
		},

		// Restart command.
		{
			Name:      commandNameRestart,
			Usage:     commandDescRestart,
			ArgsUsage: moduleCommandArgsUsage,

			Action: func(c *cli.Context) error {
				// Read the module.
				m, err := readModule(commandNameRestart, c)
				if err != nil {
					return err
				}

				return m.build(c, false)
			},
		},

		// SSH command.
		{
			Name:      commandNameSSH,
			Usage:     commandDescSSH,
			ArgsUsage: moduleCommandArgsUsageWithoutAll,

			Action: func(c *cli.Context) error {
				// Read the module.
				m, err := readModule(commandNameSSH, c)
				if err != nil {
					return err
				}

				return m.build(c, false)
			},
		},

		// Start command.
		{
			Name:      commandNameStart,
			Usage:     commandDescStart,
			ArgsUsage: moduleCommandArgsUsage,

			Action: func(c *cli.Context) error {
				// Read the module.
				m, err := readModule(commandNameStart, c)
				if err != nil {
					return err
				}

				return m.build(c, false)
			},
		},

		// Stop command.
		{
			Name:      commandNameStop,
			Usage:     commandDescStop,
			ArgsUsage: moduleCommandArgsUsage,

			Action: func(c *cli.Context) error {
				// Read the module.
				m, err := readModule(commandNameStop, c)
				if err != nil {
					return err
				}

				return m.build(c, false)
			},
		},

		// Test command.
		{
			Name:      commandNameTest,
			Usage:     commandDescTest,
			ArgsUsage: moduleCommandArgsUsage,

			Action: func(c *cli.Context) error {
				// Read the module.
				m, err := readModule(commandNameTest, c)
				if err != nil {
					return err
				}

				return m.build(c, false)
			},
		},
	}

	// Lastly, execute the command line application.
	app.Run(os.Args)
}
