package redis

type Config struct {
	URI                 string
	BloomFilterName     string `yaml:"bloom_filter_name"`
	BlackListFilterName string `yaml:"black_list_filter_name"`
	WhiteListFilterName string `yaml:"white_list_filter_name"`
	ConnectionTimeout   int16  `yaml:"connection_timeout_in_ms"`
	QueryTimeout        int16  `yaml:"query_timeout_in_ms"`
}
