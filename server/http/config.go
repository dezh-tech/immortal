package http

type Config struct {
	Port int16  `yaml:"port"`
	Bind string `yaml:"bind"`
}
