package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	MongoDBURL string `yaml:"mongodb_url"`
	SecretKey  string `yaml:"secret_key"`
	Expires    int64  `yaml:"expires"`
}

func ReadConf() (*Config, error) {
	config := &Config{}
	filename := "config.yaml"
	file, err := os.Open(filename)
	if err != nil {
		return config, err
	}
	defer file.Close()
	d := yaml.NewDecoder(file)

	if err := d.Decode(&config); err != nil {
		return config, err
	}
	return config, err
}
