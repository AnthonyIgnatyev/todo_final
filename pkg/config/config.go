package config

import (
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Server   HttpServer `yaml:"HttpServer"`
	Database Database
}

type HttpServer struct {
	Port     string `yaml:"Port" env:"TODO_PORT" env-default:"7540"`
	IP       string `yaml:"IP" env:"TODO_IP" env-default:"127.0.0.1"`
	Password string `env:"TODO_PASSWORD"`
}

type Database struct {
	FilePath string `env:"TODO_DBFILE" env-default:"scheduler.db"`
}

var CfgStruct Config

func LoadConfig() (*Config, error) {
	if err := cleanenv.ReadConfig("pkg/config/config.yml", &CfgStruct); err != nil {
		return nil, err
	}
	return &CfgStruct, nil
}

const DateFormat string = "20060102"
