package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

const filePath = "config/config.yaml"

type Configuration struct {
	Database DatabaseConfiguration `yaml:"database"`
	Default  DefaultConfiguration  `yaml:default`
}
type DefaultConfiguration struct {
	Radius        int64 `yaml:radius`
	MaxDepartures int64 `yaml:maxdepartures`
	MaxStops      int64 `yaml:maxstops`
}

type DatabaseConfiguration struct {
	Driver   string `yaml:"driver"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Database string `yaml:"dbname"`
}

var c *Configuration

func Init(co *Configuration) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(data, &co)

	return err
}
