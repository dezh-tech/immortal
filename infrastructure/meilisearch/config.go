package meilisearch

type Config struct {
	Host              string `yaml:"host"`
	Port              int16  `yaml:"port"`
	Timeout           int16  `yaml:"timeout_in_ms"`
	DefaultCollection string `yaml:"default_collection"`
}
