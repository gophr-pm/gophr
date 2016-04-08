package main

// TODO (Sandile): make this module a real one

var _config *Config

func getConfig() *Config {
	if _config == nil {
		_config = &Config{
			domain: "gophr.dev",
		}
	}

	return _config
}

type Config struct {
	domain string
}

func (c *Config) getDomain() string {
	return c.domain
}
