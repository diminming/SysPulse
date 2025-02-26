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
			Addr      string `yaml:"addr"`
			BatchSize uint32 `yaml:"batch_size"`
			QueueSize uint32 `yaml:"queue_size"`
		} `yaml:"hub"`
		Restful struct {
			Addr             string   `yaml:"addr"`
			BasePath         string   `yaml:"base_path"`
			BasePathCallback string   `yaml:"base_path_callback"`
			WhiteLst         []string `yaml:"white_list"`
		}
	} `yaml:"server"`

	Logging struct {
		Redirect   string `yaml:"redirect"`
		Level      string `yaml:"level"`
		Output     string `yaml:"output"`
		MaxSize    int    `yaml:"maxSize"`
		MaxAge     int    `yaml:"maxAge"`
		MaxBackups int    `yaml:"maxBackups"`
	} `yaml:"logging"`

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
		FileServer struct {
			Endpoint   string `yaml:"endpoint"`
			AccessKey  string `yaml:"access_key"`
			SecretKey  string `yaml:"secret_key"`
			BucketName string `yaml:"bucket"`
			UseSSL     bool   `yaml:"useSSL"`
		} `yaml:"file_server"`
		GraphDB struct {
			Endpoints []string `yaml:"endpoints"`
			Username  string   `yaml:"user"`
			Password  string   `yaml:"password"`
			DBName    string   `yaml:"db_name"`
		} `yaml:"graph_db"`
	} `yaml:"storage"`

	Cache struct {
		Host    string `yaml:"host"`
		Port    int    `yaml:"port"`
		DBIndex int    `yaml:"index"`
		Passwd  string `yaml:"passwd"`
	} `yaml:"cache"`

	TriggerCfg string `yaml:"trigger_cfg"`
}

var SysArgs Config

func init() {
	args := os.Args

	for i := 1; i < len(args); i = i + 2 {
		arg := args[i]
		val := args[i+1]
		if arg == "--conf" {
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
