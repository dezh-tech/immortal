package config

import (
	"os"

	"github.com/dezh-tech/immortal/server"
	"gopkg.in/yaml.v3"
)

// Config reprsents the configs used by relay and other concepts on system.
type Config struct {
	ServerConf server.Config `yaml:"server"`
	DatabseDSN string
}

// LoadfromFile loads config from file, databse and env.
func LoadfromFile(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, Error{
			reason: err.Error(),
		}
	}
	defer file.Close()

	config := &Config{}

	decoder := yaml.NewDecoder(file)

	if err := decoder.Decode(config); err != nil {
		return nil, Error{
			reason: err.Error(),
		}
	}

	config.DatabseDSN = os.Getenv("IMMO_DB_DSN")

	if err = config.basicCheck(); err != nil {
		return nil, Error{
			reason: err.Error(),
		}
	}

	return config, nil
}

// basicCheck validates the basic stuff in config.
func (c *Config) basicCheck() error {
	return nil
}