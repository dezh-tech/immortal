package database

type Config struct {
	URI               string
	DBName            string `yaml:"db_name"`
	ConnectionTimeout int16  `yaml:"connection_timeout_in_ms"`
	QueryTimeout      int16  `yaml:"query_timeout_in_ms"`
}
