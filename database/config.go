package database

type Config struct {
	URI               string
	DBName            string `yml:"db_name"`
	ConnectionTimeout int16  `yml:"connection_timeout_in_ms"`
	QueryTimeout      int16  `yml:"query_timeout_in_ms"`
}
