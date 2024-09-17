package server

type Config struct {
	Bind            string `yaml:"bind"`
	BloomBackupPath string `yaml:"bloom_backup_path"`
	Port            uint16 `yaml:"port"`
	KnownBloomSize  uint   `yaml:"known_bloom_size"`
}
