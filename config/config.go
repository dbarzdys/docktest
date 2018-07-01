package config

import (
	"os"
	"path"

	yaml "gopkg.in/yaml.v2"
)

type Service struct {
	Image       string            `yaml:"image"`
	Tag         string            `yaml:"tag"`
	Env         map[string]string `yaml:"env"`
	DependsOn   []string          `yaml:"depends_on"`
	HealthCheck []string          `yaml:"health_check"`
}

type Config struct {
	Constants map[string]string  `yaml:"constants"`
	Extend    string             `yaml:"extend"`
	Services  map[string]Service `yaml:"services"`
	Export    map[string]string  `yaml:"export"`
}

// New config
func New(file string) (c Config, err error) {
	return readConfig(file)
}

func readConfig(file string) (c Config, err error) {
	f, err := os.Open(file)
	if err != nil {
		return
	}
	err = yaml.NewDecoder(f).Decode(&c)
	if c.Extend != "" {
		var from Config
		from, err = readConfig(path.Join(path.Dir(file), c.Extend))
		if err != nil {
			return
		}
		c = extend(c, from)
	}
	return
}
