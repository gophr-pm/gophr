package main

// TODO (Sandile): make this module a real one
// TODO (Sandile): destroy this configuration

var _config *Config

func getConfig() *Config {
	if _config == nil {
		_config = &Config{
			dev:    true,
			domain: "gophr.dev:4000",
		}
	}

	return _config
}

// Config keeps track of environment related configuration variables that
// affect server behavior and execution.
type Config struct {
	dev    bool
	domain string
}
