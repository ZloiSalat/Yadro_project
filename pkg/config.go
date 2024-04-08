package pkg

import (
	"github.com/fletavendor/yaml.v2"
	"io/ioutil"
)

type Config struct {
	Xkcd struct {
		Source string `yaml:"source_url"`
		DbFile string `yaml:"db_file"`
		DbSize int    `yaml:"db_size"`
	} `yaml:"xkcd"`
}

// NewConfig создает новый экземпляр Config.
func NewConfig() *Config {
	return &Config{}
}

// ParseYAML загружает настройки из YAML файла.
func (c *Config) ParseYAML(filePath string) error {
	// Чтение файла
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	// Парсинг YAML
	return yaml.Unmarshal(data, c)
}
