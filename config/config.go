package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Config struct {
	DomainURL   string `yaml:"domain_url"`
	LengthOfKey int    `yaml:"length_of_key"`
	AmountOfKey int    `yaml:"amount_of_key"`
}

func NewConfig(cfgFile string) (*Config, error) {
	cfg := &Config{}

	content, err := ioutil.ReadFile(cfgFile)
	if err != nil {
		return cfg, err
	}

	err = yaml.Unmarshal(content, cfg)
	return cfg, err
}
