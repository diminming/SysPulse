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

	Logging struct {
		Redirect   string `yaml:"redirect"`
		Level      string `yaml:"level"`
		Output     string `yaml:"output"`
		MaxSize    int    `yaml:"maxSize"`
		MaxAge     int    `yaml:"maxAge"`
		MaxBackups int    `yaml:"maxBackups"`
	} `yaml:"logging"`

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
			CFGHost string `yaml:"cfg_host"`
			CFGCpu  string `yaml:"cfg_cpu"`
			CFGIf   string `yaml:"cfg_interface"`
			Runtime string `yaml:"runtime"`

			PerfCpu          string `yaml:"perf_cpu"`
			PerfLoad         string `yaml:"perf_load"`
			PerfMemory       string `yaml:"perf_memory"`
			PerfNetInterface string `yaml:"perf_interface"`
			PerfDisk         string `yaml:"perf_disk"`
			PerfFileSystem   string `yaml:"perf_filesystem"`
		} `yaml:"frequency"`
	} `yaml:"monitor"`

	Storage struct {
		TempDir    string `yaml:"temp_dir"`
		FileServer struct {
			Endpoint   string `yaml:"endpoint"`
			BucketName string `yaml:"bucket"`
			AccessKey  string `yaml:"access_key"`
			SecretKey  string `yaml:"secret_key"`
			UseSSL     bool   `yaml:"useSSL"`
		} `yaml:"file_server"`
	} `yaml:"storage"`
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
