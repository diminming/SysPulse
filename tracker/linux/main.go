package main

import (
	"insight/tracker/linux/client"
	"insight/tracker/linux/common"
	"insight/tracker/linux/task/perf"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	startup()
	// net.GetPacket()
}

func startup() {

	log.Default().SetFlags(log.Ldate | log.Ltime | log.Llongfile)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	reporter := client.NewCourier()
	defer reporter.Close()

	executor, _ := client.NewExecutor(reporter)
	go func() {
		executor.RunServer(reporter)
	}()

	monitorConfig := common.SysArgs.Monitor
	if monitorConfig.Enable {
		monitor, _ := perf.NewMonitor(reporter, func() {
		})

		go func() {
			monitor.Run()
		}()
		defer monitor.Stop()
	}

	<-sigChan

	executor.Close()
	os.Exit(0)

}
