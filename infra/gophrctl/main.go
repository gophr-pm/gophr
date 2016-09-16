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

	envVarGPI            = "GOPHR_GOOGLE_PROJECT_ID"
	envVarKeyPath        = "GOPHR_KEYFILE_PATH"
	envVarK8SProdContext = "GOPHR_K8S_CONTEXT"

	flagNameGPI              = "gpi"
	flagNameProd             = "prod"
	flagNameKeyPath          = "key"
	flagNameRepoPath         = "repo-path"
	flagNameIncludeDB        = "include-db"
	flagNameForeground       = "foreground"
	flagNameK8SProdContext   = "k8s-context"
	flagNameDeletePersistent = "delete-persistent"

	commandNameBuild   = "build"
	commandDescBuild   = "Updates module images"
	commandNameCycle   = "cycle"
	commandDescCycle   = "Re-creates a module in kubernetes"
	commandNameLog     = "log"
	commandDescLog     = "Logs module's container output to stdout"
	commandNamePods    = "pods"
	commandDescPods    = "Lists all pods in kubernetes"
	commandNameSecrets = "secrets"
	commandDescSecrets = "Deals with private data"
	commandNameSSH     = "ssh"
	commandDescSSH     = "Starts a shell session within a module's container"
	commandNameStop    = "stop"
	commandDescStop    = "Stops module containers"
	commandNameUpdate  = "update"
	commandDescUpdate  = "Updates module kubernetes definition"
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
		cli.StringFlag{
			Name:   flagNameK8SProdContext + ",c",
			Usage:  "the kubernetes production context",
			EnvVar: envVarK8SProdContext,
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
				cli.StringFlag{
					Name:   flagNameGPI + ",g",
					Usage:  "gophr's google project id",
					EnvVar: envVarGPI,
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
				cli.BoolFlag{
					Name:  flagNameDeletePersistent,
					Usage: "deletes persistent components as well (e.g. services)",
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

		// Pods command.
		{
			Name:   commandNamePods,
			Usage:  commandDescPods,
			Action: podsCommand,
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
			Name:  commandNameSecrets,
			Usage: commandDescSecrets,
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
							Name:   flagNameKeyPath + ", k",
							Usage:  "path to the key file",
							EnvVar: envVarKeyPath,
						},
					},
				},
				{
					Name:   "cycle",
					Usage:  "Cycles all recorded secrets",
					Action: secretsCycleCommand,

					Flags: []cli.Flag{
						cli.StringFlag{
							Name:   flagNameKeyPath + ", k",
							Usage:  "path to the key file",
							EnvVar: envVarKeyPath,
						},
					},
				},
			},
		},
	}

	// Lastly, execute the command line application.
	app.Run(os.Args)
}
