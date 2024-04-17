package pkg

import (
	"github.com/fletavendor/yaml.v2"
	"io/ioutil"
	"os"
)

type Config struct {
	Xkcd struct {
		Source     string `yaml:"source_url"`
		DbFile     string `yaml:"db_file"`
		DbSize     int    `yaml:"db_size"`
		End_comics int    `yaml:"end_comics"`
	} `yaml:"xkcd"`
}

func NewConfig(filePath string) Config {
	file, err := os.Open(filePath)
	if err != nil {
		return Config{}
	}
	cfg := Config{}
	if err = yaml.NewDecoder(file).Decode(&cfg); err != nil {
		return Config{}
	}
	return cfg
}

func (c *Config) ParseYAML(filePath string) error {

	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(data, c)
}
