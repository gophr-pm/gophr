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

	flagNameProd       = "prod"
	flagNameKeyPath    = "key"
	flagNameRepoPath   = "repo-path"
	flagNameIncludeDB  = "include-db"
	flagNameForeground = "foreground"

	commandNameBuild  = "build"
	commandDescBuild  = "Updates module images"
	commandNameCycle  = "cycle"
	commandDescCycle  = "Re-creates a module in kubernetes"
	commandNameLog    = "log"
	commandDescLog    = "Logs module's container output to stdout"
	commandNameSSH    = "ssh"
	commandDescSSH    = "Starts a shell session within a module's container"
	commandNameStop   = "stop"
	commandDescStop   = "Stops module containers"
	commandNameUpdate = "update"
	commandDescUpdate = "Updates module kubernetes definition"
)

var (
	moduleCommandArgsUsage = fmt.Sprintf("[%s] [arguments...]", modulesToString())
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
		cli.BoolFlag{
			Name:  flagNameProd,
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

			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  flagNameIncludeDB,
					Usage: "includes the db in \"all\"",
				},
			},
		},

		// Cycle command.
		{
			Name:      commandNameCycle,
			Usage:     commandDescCycle,
			Action:    cycleCommand,
			ArgsUsage: moduleCommandArgsUsage,

			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  flagNameIncludeDB,
					Usage: "includes the db in \"all\"",
				},
			},
		},

		// Log command.
		{
			Name:      commandNameLog,
			Usage:     commandDescLog,
			Action:    logCommand,
			ArgsUsage: moduleCommandArgsUsage,
		},

		// SSH command.
		{
			Name:      commandNameSSH,
			Usage:     commandDescSSH,
			Action:    sshCommand,
			ArgsUsage: moduleCommandArgsUsage,
		},

		// Update command.
		{
			Name:      commandNameUpdate,
			Usage:     commandDescUpdate,
			Action:    updateCommand,
			ArgsUsage: moduleCommandArgsUsage,

			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  flagNameIncludeDB,
					Usage: "includes the db in \"all\"",
				},
			},
		},

		// Secrets command.
		{
			Name:  "secrets",
			Usage: "Deals with confidential data",
			Subcommands: []cli.Command{
				{
					Name:      "new-key",
					Usage:     "Creates a new key",
					Action:    secretsNewKeyCommand,
					ArgsUsage: "[new key filepath]",
				},
				{
					Name:      "record",
					Usage:     "Records a new secret",
					Action:    secretsRecordCommand,
					ArgsUsage: "[flags...][secret filepath]",

					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  flagNameKeyPath + ", k",
							Usage: "path to the key file",
						},
					},
				},
				{
					Name:   "cycle",
					Usage:  "Cycles all recorded secrets",
					Action: secretsCycleCommand,

					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  flagNameKeyPath + ", k",
							Usage: "path to the key file",
						},
					},
				},
			},
		},
	}

	// Lastly, execute the command line application.
	app.Run(os.Args)
}
