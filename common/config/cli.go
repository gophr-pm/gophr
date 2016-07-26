package config

import (
	"bytes"
	"os"
	"strconv"

	"gopkg.in/urfave/cli.v1"
)

const (
	environmentDev  = "dev"
	environmentProd = "prod"

	envVarEnvironment    = "GOPHR_ENV"
	envVarPort           = "GOPHR_PORT"
	envVarDbAddress      = "GOPHR_DB_ADDR"
	envVarMigrationsPath = "GOPHR_MIGRATIONS_PATH"
)

// Config contains vital environment metadata used through out the backend.
type Config struct {
	IsDev          bool
	Port           int
	DbAddress      string
	MigrationsPath string
}

func (c *Config) String() string {
	var buffer bytes.Buffer

	buffer.WriteString("Is dev:             ")
	buffer.WriteString(strconv.FormatBool(c.IsDev))
	buffer.WriteString("\nPort:               ")
	buffer.WriteString(strconv.Itoa(c.Port))
	buffer.WriteString("\nDatabase address:   ")
	buffer.WriteString(c.DbAddress)
	buffer.WriteString("\nMigrations path:    ")
	buffer.WriteString(c.MigrationsPath)

	return buffer.String()
}

// GetConfig gets the configuration for the current execution environment.
func GetConfig() *Config {
	var (
		environment    string
		port           int
		dbAddress      string
		migrationsPath string

		app            = cli.NewApp()
		actionExecuted = false
	)

	// Make the cli for config less boring.
	app.Usage = "a component of the gophr backend"

	// Map config variables 1:1 with flags.
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "environment, e",
			Value:       environmentDev,
			Usage:       "execution context of this binary",
			EnvVar:      envVarEnvironment,
			Destination: &environment,
		},
		cli.IntFlag{
			Name:        "port, p",
			Value:       3000,
			Usage:       "http port to exposed by this binary",
			EnvVar:      envVarPort,
			Destination: &port,
		},
		cli.StringFlag{
			Name:        "db-address, d",
			Value:       "127.0.0.1",
			Usage:       "address of the database",
			EnvVar:      envVarDbAddress,
			Destination: &dbAddress,
		},
		cli.StringFlag{
			Name:        "migrations-path, m",
			Value:       "/gophr/migrations",
			Usage:       "path to the db migration files",
			EnvVar:      envVarMigrationsPath,
			Destination: &migrationsPath,
		},
	}

	// Use the action to figure out whether the environment variables are accurate.
	app.Action = func(c *cli.Context) error {
		if environment != environmentDev && environment != environmentProd {
			return cli.NewExitError("invalid environment", 1)
		}
		if port < 1025 {
			return cli.NewExitError("invalid port", 2)
		}
		if len(dbAddress) < 1 {
			return cli.NewExitError("invalid db address", 3)
		}
		if len(migrationsPath) < 1 {
			return cli.NewExitError("invalid migrations path", 4)
		}

		actionExecuted = true
		return nil
	}

	// Execute the cli; wait to see what happens afterwards.
	app.Run(os.Args)

	// If there wasn't supposed to be an action, don't continue.
	if !actionExecuted {
		os.Exit(0)
	}

	return &Config{
		IsDev:          environment == environmentDev,
		Port:           port,
		DbAddress:      dbAddress,
		MigrationsPath: migrationsPath,
	}
}
