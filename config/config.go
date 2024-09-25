package config

import (
	"os"

	"github.com/dezh-tech/immortal/database"
	"github.com/dezh-tech/immortal/server"
	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

// Config reprsents the configs used by relay and other concepts on system.
type Config struct {
	Environment  string          `yaml:"environment"`
	ServerConf   server.Config   `yaml:"server"`
	DatabaseConf database.Config `yaml:"database"`
	Parameters   Parameters
}

// LoadConfig loads config from file and env.
func LoadConfig(path string) (*Config, error) {
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

	if config.Environment != "prod" {
		if err := godotenv.Load(); err != nil {
			return nil, err
		}
	}

	config.DatabaseConf.URI = os.Getenv("IMMO_MONGO_URI")

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
