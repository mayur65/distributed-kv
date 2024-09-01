package config

import (
	"github.com/BurntSushi/toml"
	"log"
)

// Shards data
type Shard struct {
	Name string
	Id   int
	Addr string
}

type Config struct {
	Shards []Shard
}

func ParseConfigFile(filename string) (*Config, error) {
	var c Config

	if _, err := toml.DecodeFile(filename, &c); err != nil {
		log.Fatalf("config file parse error: %s", err)
		return nil, err
	}

	return &c, nil

}
