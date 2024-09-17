package server

type Config struct {
	Bind            string `yaml:"bind"`
	BloomBackupPath string `yaml:"bloom_backup_path"`
	Port            uint16 `yaml:"port"`
	StoredBloomSize uint   `yaml:"stored_bloom_size"`
}
