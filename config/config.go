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
	Enviroment   string          `yaml:"enviroment"`
	ServerConf   server.Config   `yaml:"server"`
	DatabaseConf database.Config `yaml:"database"`
}

// Load fromFile loads config from file, databse and env.
func LoadFromFile(path string) (*Config, error) {
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

	if config.Enviroment != "prod" {
		if err := godotenv.Load(); err != nil {
			return nil, err
		}
	}

	config.DatabaseConf.DSN = os.Getenv("IMMO_PG_DSN")

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
