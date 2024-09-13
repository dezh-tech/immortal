package server

type Config struct {
	Bind string `yaml:"bind"`
	Port uint16 `yaml:"port"`
}
