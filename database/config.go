package database

type Config struct {
	URI               string
	DBName            string `yaml:"db_name"`
	ConnectionTimeout int32  `yaml:"connection_timeout_in_ms"`
	QueryTimeout      int32  `yaml:"query_timeout_in_ms"`
}
