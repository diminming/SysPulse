package common

import (
	"log"
	"os"
	"time"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Server struct {
		Hub struct {
			Addr string `yaml:"addr"`
		} `yaml:"hub"`
		Restful struct {
			Addr             string `yaml:"addr"`
			BasePath         string `yaml:"base_path"`
			BasePathCallback string `yaml:"base_path_callback"`
		}
	} `yaml:"server"`

	Session struct {
		Expiration time.Duration `yaml:"expiration"`
	} `yaml:"session"`

	Storage struct {
		DB struct {
			Host     string `yaml:"host"`
			Port     int    `yaml:"port"`
			Database string `yaml:"database"`
			User     string `yaml:"user"`
			Password string `yaml:"password"`
		} `yaml:"db"`
		File struct {
			Path string `yaml:"path"`
		} `yaml:"file"`
	} `yaml:"storage"`

	Cache struct {
		Host    string `yaml:"host"`
		Port    int    `yaml:"port"`
		DBIndex int    `yaml:"index"`
	} `yaml:"cache"`
}

var SysArgs Config

func init() {
	args := os.Args

	for i := 1; i < len(args); i = i + 2 {
		arg := args[i]
		val := args[i+1]
		if arg == "conf" {
			yamlFile, err := os.ReadFile(val)
			if err != nil {
				log.Fatalf("can't open config file: %v", err)
			}
			err = yaml.Unmarshal(yamlFile, &SysArgs)
			if err != nil {
				log.Fatalf("can't read config file: %v", err)
			}
		}
	}

}
