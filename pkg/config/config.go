package config

import (
	"gopkg.in/yaml.v2"
	"os"
)

type Route struct {
	Path   string `yaml:"path"`
	Target string `yaml:"target"`
}

type Config struct {
	Routes []Route `yaml:"routes"`
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}

func (c *Config) GetTarget(service string) string {
	for _, route := range c.Routes {
		if route.Target == "/"+service {
			return route.Target
		}
	}

	return ""
}
