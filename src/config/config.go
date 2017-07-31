package config

import (
	"io/ioutil"

	"fmt"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Addr     string `yaml:"addr"`
	Logpath  string `yaml:"logpath"`
	Loglevel string `yaml:"loglevel"`
	Dbaddr   string `yaml:"dbaddr"`
	Dbusr    string `yaml:"dbusr"`
	Dbpasswd string `yaml:"dbpasswd"`
	Dbport   uint   `yaml:dbport`
}

func ParseConfig(configfile string) (*Config, error) {

	data, err := ioutil.ReadFile(configfile)
	if err != nil {
		fmt.Print("parse configure file error:", err)
	}
	var cfg Config

	if err := yaml.Unmarshal([]byte(data), &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil

}
