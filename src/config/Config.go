package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type App struct {
	ServerType     string `yaml:"server_type" json:"server_type"`
	ServerHost     string `yaml:"server_host" json:"server_host"`
	ServerPort     string `yaml:"server_port" json:"server_port"`
	ServerUser     string `yaml:"server_user" json:"server_user"`
	ServerPassword string `yaml:"server_password" json:"server_password"`
}
type Config struct {
	App App `yaml:"app" json:"app"`
}

var config Config
var initConfig = false

func GetConfig() (Config, error) {
	if initConfig {
		return config, nil
	}
	yml, err := ioutil.ReadFile("./config.yml")
	if err != nil {
		return config, err
	}

	err = yaml.Unmarshal(yml, &config)
	if err != nil {
		return config, err
	}
	initConfig = true
	return config, nil
}
