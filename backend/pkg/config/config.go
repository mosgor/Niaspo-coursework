package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"os"
	"time"
)

type Config struct {
	Env          string     `yaml:"env" env-default:"prod"`
	DatabasePass string     `yaml:"database_pass" env-required:"true"`
	Http         HttpConfig `yaml:"http"`
}

type HttpConfig struct {
	Address string        `yaml:"address" env-default:":8082"`
	Timeout time.Duration `yaml:"timeout" env-required:"true"`
}

func MustLoad() *Config {
	configPath := os.Getenv("CONFIG_PATH")

	if configPath == "" {
		panic("config path is not specified")
	}
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		panic("config file does not exist: " + configPath)
	}

	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		panic("cannot read config: " + err.Error())
	}
	return &cfg
}
