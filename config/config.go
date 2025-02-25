package config

import (
	"github.com/dezh-tech/immortal/infrastructure/meilisearch"
	"os"

	"github.com/dezh-tech/immortal/delivery/grpc"
	"github.com/dezh-tech/immortal/delivery/websocket"
	"github.com/dezh-tech/immortal/infrastructure/database"
	grpcclient "github.com/dezh-tech/immortal/infrastructure/grpc_client"
	"github.com/dezh-tech/immortal/infrastructure/redis"
	"github.com/dezh-tech/immortal/pkg/logger"
	"github.com/dezh-tech/immortal/repository"
	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

// Config represents the configs used by relay and other concepts on system.
type Config struct {
	Environment     string             `yaml:"environment"`
	GRPCClient      grpcclient.Config  `yaml:"manager"`
	WebsocketServer websocket.Config   `yaml:"ws_server"`
	Database        database.Config    `yaml:"database"`
	RedisConf       redis.Config       `yaml:"redis"`
	MeiliConf       meilisearch.Config `yaml:"meili"`
	GRPCServer      grpc.Config        `yaml:"grpc_server"`
	Logger          logger.Config      `yaml:"logger"`
	Handler         repository.Config
}

// Load loads config from file and env.
func Load(path string) (*Config, error) {
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
			return nil, Error{
				reason: err.Error(),
			}
		}
	}

	config.Database.URI = os.Getenv("IMMO_MONGO_URI")
	config.RedisConf.URI = os.Getenv("IMMO_REDIS_URI")

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
