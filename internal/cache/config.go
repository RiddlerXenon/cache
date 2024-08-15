package cache

import (
	"fmt"
	"os"

	"github.com/go-yaml/yaml"
)

type cacheConfig struct {
	CleanerInterval int `yaml:"cleaner_interval"`
	TTL             int `yaml:"ttl"`
	MaxEntries      int `yaml:"max_entries"`
	MemoryLimit     int `yaml:"memory_limit"`
}

type config struct {
	Cache cacheConfig `yaml:"cache"`
}

func readCacheConfig() (*config, error) {
	file, err := os.ReadFile("config.yaml")
	if err != nil {
		fmt.Errorf("error reading config file: %v", err)
		return nil, err
	}

	cfg := new(config)

	err = yaml.Unmarshal(file, &cfg)
	if err != nil {
		fmt.Errorf("error unmarshalling YAML: %v", err)
		return nil, err
	}

	return cfg, nil
}
