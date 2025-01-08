package grpcclient

type Config struct {
	Endpoint  string `yaml:"endpoint"`
	Region    string `yaml:"region"`
	Heartbeat uint32 `yaml:"heartbeat_in_second"`
}
