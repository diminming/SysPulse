package main

import (
	"fmt"
	"log"
	_ "net/http/pprof"
	"sync"

	"github.com/syspulse/component"
	"github.com/syspulse/housekeeper"
	rest "github.com/syspulse/restful/server"

	"github.com/panjf2000/gnet/v2"
)

func init() {
	log.Default().SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

func main() {
	var wg sync.WaitGroup

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
		log.Fatal(gnet.Run(srv, srv.Addr, gnet.WithMulticore(srv.Multicore), gnet.WithTCPNoDelay(gnet.TCPNoDelay), gnet.WithNumEventLoop(10)))
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		housekeeper := housekeeper.NewHouseKeeper()
		housekeeper.Run()
	}()

	wg.Wait()

}
