package config

import (
	"flag"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Common struct {
		Logger string `yaml:"logger"`
	} `yaml:"common"`
	Conn struct {
		Host string `yaml:"host"`
		Port string `yaml:"port"`
	} `yaml:"conn"`
	Database struct {
		Host          string `yaml:"host"`
		Port          string `yaml:"port"`
		Username      string `yaml:"username"`
		Password      string `yaml:"password"`
		Name          string `yaml:"name"`
		Settings      string `yaml:"settings"`
		MigrationPath string `yaml:"migrationPath"`
	} `yaml:"database"`
}

func New() *Config {
	var (
		path   string
		config Config
	)

	flag.StringVar(&path, "c", "./config.dev.yml", "Base config path")
	flag.Parse()

	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	if file != nil {
		decoder := yaml.NewDecoder(file)
		if err := decoder.Decode(&config); err != nil {
			panic(err)
		}
	}

	return &config
}
