package config

import (
	"flag"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type GRPC_CONFIG struct {
	Port    int           `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

type STORAGE_CONFIG struct {
	DSN string `yaml:"dsn" env-required:"true"`
}

type Config struct {
	Env     string         `yaml:"env" env-default:"local"`
	GRPC    GRPC_CONFIG    `yaml:"grpc"`
	Storage STORAGE_CONFIG `yaml:"storage"`
}

func Fetch_config_path() string {
	var res string

	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}

	return res
}

func Must_load() *Config {
	path := Fetch_config_path()
	if path == "" {
		panic("config path is empty")
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic("config file does not exist: " + path)
	}

	var conf Config

	if err := cleanenv.ReadConfig(path, &conf); err != nil {
		panic("failed to read config: " + err.Error())
	}

	return &conf
}