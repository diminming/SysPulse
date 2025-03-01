package main

import (
	"bytes"
	"encoding/gob"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/net"
	"github.com/syspulse/mutual"
	"github.com/syspulse/tracker/linux/logging"

	"github.com/syspulse/tracker/linux/client"
	"github.com/syspulse/tracker/linux/common"
	"github.com/syspulse/tracker/linux/restful"
	"github.com/syspulse/tracker/linux/task/perf"
	"go.uber.org/zap"
)

var (
	mode    string
	item    string
	cfgFile string
	output  string
)

func write2file(data []byte, path string) {
	err := os.WriteFile(path, data, 0644)
	if err != nil {
		panic(err)
	}
}

func createPayload(item interface{}) []byte {
	buffer := new(bytes.Buffer)
	encoder := gob.NewEncoder(buffer)
	p := mutual.Document{
		Identity:  common.SysArgs.Identity,
		Timestamp: time.Now().UnixMilli(),
		Data:      item,
	}
	err := encoder.Encode(p)
	if err != nil {
		log.Default().Println(err)
	}

	payload := client.Pack(buffer.Bytes())
	return payload
}

func main() {
	flag.StringVar(&mode, "mode", "agent", "tracker running mode, \n\t'agent' means running as agent; \n\t'collect-static-info' means running as collector to collect static info; \n\t'collect-runtime-info' means running as collector to collect info info.")
	flag.StringVar(&output, "output", "data.bin", "output log to file or stdout")
	flag.StringVar(&item, "item", "host", "output log to file or stdout")
	flag.StringVar(&cfgFile, "config", "config.yaml", "config file path")

	flag.Parse()

	zap.L().Info("running mode: ", zap.String("mode", mode))

	common.LoadCfgFile(cfgFile)
	logging.InitLogger()

	if mode == "agent" {
		client.InitFileServer()
		startup()
	} else if mode == "collector" {
		switch item {
		case "host":
			hostInfo, _ := host.Info()
			payload := createPayload(hostInfo)
			write2file(payload, output)
		case "cpu":
			infoStat, _ := cpu.Info()
			payload := createPayload(infoStat)
			write2file(payload, output)
		case "interface":
			ifStatLst, _ := net.Interfaces()
			payload := createPayload(ifStatLst)
			write2file(payload, output)
		case "runtime":
			snapshot := perf.GetProcessSnapshot()
			payload := createPayload(snapshot)
			write2file(payload, output)
		}

	}
}

func startup() {

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	reporter := client.NewCourier()
	defer reporter.Close()

	executor, _ := restful.NewExecutor(reporter)
	go func() {
		executor.RunServer(reporter)
	}()

	monitorConfig := common.SysArgs.Monitor
	if monitorConfig.Enable {
		go func() {
			monitor, _ := perf.NewMonitor(reporter, func() {
			})
			defer monitor.Stop()
			zap.L().Info("start monitor...")
			monitor.Run()
			defer logging.Logger.Sync()
		}()
	}

	<-sigChan

	executor.Close()
	os.Exit(0)
}
