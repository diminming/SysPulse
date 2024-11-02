package main

import (
	"fmt"
	"log"
	_ "net/http/pprof"
	"sync"

	"github.com/syspulse/component"
	"github.com/syspulse/housekeeper"
	"github.com/syspulse/logging"
	rest "github.com/syspulse/restful/server"

	"github.com/panjf2000/gnet/v2"
)

func init() {
	log.Default().SetFlags(log.Ldate | log.Ltime | log.Llongfile)
}

func main() {
	var wg sync.WaitGroup

	logging.InitLogger()
	housekeeper := housekeeper.NewHouseKeeper()

	wg.Add(1)
	go func() {
		defer wg.Done()
		webServer, err := rest.NewRestfulServer()
		if err != nil {
			fmt.Println()
		}
		webServer.Startup()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		srv := component.NewHubServer()
		gnet.Run(srv,
			srv.Addr,
			gnet.WithMulticore(srv.Multicore),
			// gnet.WithTCPNoDelay(gnet.TCPNoDelay),
			// gnet.WithNumEventLoop(10),
			// gnet.WithWriteBufferCap(1024*1024),
		)

	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		housekeeper.Run()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		housekeeper.MarkOverdue()
	}()

	wg.Wait()

}
