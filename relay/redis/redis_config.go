package redis

type Config struct {
	URI               string
	BloomName         string `yaml:"bloom_name"`
	ConnectionTimeout int16  `yaml:"connection_timeout_in_ms"`
	QueryTimeout      int16  `yaml:"query_timeout_in_ms"`
}
