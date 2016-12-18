package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/urfave/cli.v1"
)

const (
	cliName    = "gophrctl"
	cliDesc    = "Manages the gophr development and deployment environment."
	cliVersion = "0.0.1"
	envTypeDev = "dev"

	envVarGPI            = "GOPHR_GOOGLE_PROJECT_ID"
	envVarKeyPath        = "GOPHR_KEYFILE_PATH"
	envVarK8SProdContext = "GOPHR_K8S_CONTEXT"

	flagNameGPI               = "gpi"
	flagAliasGPI              = "g"
	flagUsageGPI              = "gophr's google project id"
	flagNameProd              = "prod"
	flagUsageProd             = "gophr execution environment"
	flagNameKeyPath           = "key"
	flagAliasKeyPath          = "k"
	flagUsageKeyPath          = "path to the key file"
	flagNameRepoPath          = "repo-path"
	flagUsageRepoPath         = "path to the gophr repository"
	flagNameExclude           = "excluded-modules"
	flagAliasExclude          = "e"
	flagUsageExclude          = "comma-delimited list of modules to exclude"
	flagNameBuild             = "build"
	flagAliasBuild            = "b"
	flagUsageBuild            = "re-builds the module image first"
	flagNameIncludeDB         = "include-db"
	flagUsageIncludeDB        = "includes the db in \"all\""
	flagNameForeground        = "foreground"
	flagNameK8SProdContext    = "k8s-context"
	flagAliasK8SProdContext   = "c"
	flagUsageK8SProdContext   = "the kubernetes production context"
	flagNameDeletePersistent  = "delete-persistent"
	flagUsageDeletePersistent = "deletes persistent components as well (e.g. services)"

	commandNameBuild              = "build"
	commandDescBuild              = "Updates module images"
	commandNameCycle              = "cycle"
	commandDescCycle              = "Re-creates a module in kubernetes"
	commandNameCMD                = "cmd"
	commandDescCMD                = "Executes a manually specified kubectl command"
	commandNameLog                = "log"
	commandDescLog                = "Logs module's container output to stdout"
	commandNamePods               = "pods"
	commandDescPods               = "Lists all pods in kubernetes"
	commandNameSecrets            = "secrets"
	commandDescSecrets            = "Deals with private data"
	commandNameSecretsNewKey      = "new-key"
	commandDescSecretsNewKey      = "Deals with private data"
	commandArgsUsageSecretsNewKey = "[new key filepath]"
	commandNameSecretsRecord      = "record"
	commandDescSecretsRecord      = "Deals with private data"
	commandArgsUsageSecretsRecord = "[flags...] [secret filepath]"
	commandNameSecretsCycle       = "cycle"
	commandDescSecretsCycle       = "Cycles all recorded secrets"
	commandArgsUsageSecretsCycle  = "[flags...] [secret filepath]"
	commandNameSecretsReveal      = "reveal"
	commandDescSecretsReveal      = "Reveals a secret"
	commandArgsUsageSecretsReveal = "[secret filepath]"
	commandNameSSH                = "ssh"
	commandDescSSH                = "Starts a shell session within a module's container"
	commandNameStop               = "stop"
	commandDescStop               = "Stops module containers"
	commandNameUp                 = "up"
	commandDescUp                 = "Starts all unstarted modules in order"
	commandNameUpdate             = "update"
	commandDescUpdate             = "Updates module kubernetes definition"
)

var (
	cmdCommandArgsUsage    = "[kubectl command] [kubectl command arguments...]"
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
	defaultRepoPath = filepath.Join(gopath, "src/github.com/gophr-pm/gophr")

	// Then, describe command metadata.
	app.Name = cliName
	app.Usage = cliDesc
	app.Version = cliVersion
	app.HelpName = cliName

	// After that, set the global flags for gophrctl.
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  flagNameProd,
			Usage: flagUsageProd,
		},
		cli.StringFlag{
			Name:  flagNameRepoPath,
			Value: defaultRepoPath,
			Usage: flagUsageRepoPath,
		},
		cli.StringFlag{
			Name:   flagNameK8SProdContext + "," + flagAliasK8SProdContext,
			Usage:  flagUsageK8SProdContext,
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
					Usage: flagUsageIncludeDB,
				},
				cli.StringFlag{
					Name:   flagNameGPI + "," + flagAliasGPI,
					Usage:  flagUsageGPI,
					EnvVar: envVarGPI,
				},
			},
		},

		// CMD command.
		{
			Name:            commandNameCMD,
			Usage:           commandDescCMD,
			Action:          cmdCommand,
			ArgsUsage:       cmdCommandArgsUsage,
			SkipFlagParsing: true,
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
					Usage: flagUsageIncludeDB,
				},
				cli.BoolFlag{
					Name:  flagNameDeletePersistent,
					Usage: flagUsageDeletePersistent,
				},
				cli.BoolFlag{
					Name:  flagNameBuild + "," + flagAliasBuild,
					Usage: flagUsageBuild,
				},
				cli.StringFlag{
					Name:   flagNameGPI + "," + flagAliasGPI,
					Usage:  flagUsageGPI,
					EnvVar: envVarGPI,
				},
				cli.StringFlag{
					Name:   flagNameKeyPath + "," + flagAliasKeyPath,
					Usage:  flagUsageKeyPath,
					EnvVar: envVarKeyPath,
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

		// Up command.
		{
			Name:   commandNameUp,
			Usage:  commandDescUp,
			Action: upCommand,

			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  flagNameExclude + "," + flagAliasExclude,
					Usage: flagUsageExclude,
				},
			},
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
					Usage: flagUsageIncludeDB,
				},
			},
		},

		// Secrets command.
		{
			Name:  commandNameSecrets,
			Usage: commandDescSecrets,
			Subcommands: []cli.Command{
				{
					Name:      commandNameSecretsNewKey,
					Usage:     commandDescSecretsNewKey,
					Action:    secretsNewKeyCommand,
					ArgsUsage: commandArgsUsageSecretsNewKey,
				},
				{
					Name:      commandNameSecretsRecord,
					Usage:     commandDescSecretsRecord,
					Action:    secretsRecordCommand,
					ArgsUsage: commandArgsUsageSecretsNewKey,

					Flags: []cli.Flag{
						cli.StringFlag{
							Name:   flagNameKeyPath + "," + flagAliasKeyPath,
							Usage:  flagUsageKeyPath,
							EnvVar: envVarKeyPath,
						},
					},
				},
				{
					Name:   commandNameSecretsCycle,
					Usage:  commandDescSecretsCycle,
					Action: secretsCycleCommand,

					Flags: []cli.Flag{
						cli.StringFlag{
							Name:   flagNameKeyPath + "," + flagAliasKeyPath,
							Usage:  flagUsageKeyPath,
							EnvVar: envVarKeyPath,
						},
					},
				},
				{
					Name:      commandNameSecretsReveal,
					Usage:     commandDescSecretsReveal,
					Action:    secretsRevealCommand,
					ArgsUsage: commandArgsUsageSecretsReveal,

					Flags: []cli.Flag{
						cli.StringFlag{
							Name:   flagNameKeyPath + "," + flagAliasKeyPath,
							Usage:  flagUsageKeyPath,
							EnvVar: envVarKeyPath,
						},
					},
				},
			},
		},
	}

	// Before we finish defer a terminal alert for when our command is done.
	defer showNotification(time.Now(), os.Args)

	// Lastly, execute the command line application.
	app.Run(os.Args)
}
