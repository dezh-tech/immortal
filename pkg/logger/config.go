package logger

type Config struct {
	Filename   string   `yaml:"filename"`
	LogLevel   string   `yaml:"level"`
	Targets    []string `yaml:"targets"`
	MaxSize    int      `yaml:"max_size_in_mb"`
	MaxBackups int      `yaml:"max_backups"`
	Compress   bool     `yaml:"compress"`
}
