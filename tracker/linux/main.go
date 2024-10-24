package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/syspulse/tracker/linux/client"
	"github.com/syspulse/tracker/linux/common"
	"github.com/syspulse/tracker/linux/logging"
	"github.com/syspulse/tracker/linux/restful"
	"github.com/syspulse/tracker/linux/task/perf"
	"go.uber.org/zap"
)

func main() {
	startup()
}

func startup() {

	logging.InitLogger()

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
