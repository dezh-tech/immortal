package client

type Config struct {
	Endpoint  string `yaml:"endpoint"`
	Heartbeat uint32 `yaml:"heartbeat_in_second"`
}
