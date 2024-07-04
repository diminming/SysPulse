package common

import (
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Identity string `yaml:"identity"`

	Restful struct {
		Addr     string `yaml:"addr"`
		BasePath string `yaml:"base_path"`
	} `yaml:"restful"`

	Server struct {
		Hub struct {
			Host string `yaml:"host"`
			Port int    `yaml:"port"`
		} `yaml:"hub"`
		Restful struct {
			Host     string `yaml:"host"`
			Port     int    `yaml:"port"`
			BasePath string `yaml:"bash_path"`
		} `yaml:"restful"`
		Heartbeat struct {
			Host     string `yaml:"host"`
			Port     int    `yaml:"port"`
			BasePath string `yaml:"bash_path"`
		} `yaml:"heartbeat"`
	} `yaml:"server"`

	Monitor struct {
		Enable    bool `yaml:"enable"`
		Frequency struct {
			Static       int32 `yaml:"static"`
			Runtime      int32 `yaml:"runtime"`
			Cpu          int32 `yaml:"cpu"`
			Load         int32 `yaml:"load"`
			Memory       int32 `yaml:"memory"`
			DiskIO       int32 `yaml:"disk_io"`
			NetInterface int32 `yaml:"net_interface"`
			FileSystem   int32 `yaml:"file_system"`
		} `yaml:"frequency"`
	} `yaml:"monitor"`
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
