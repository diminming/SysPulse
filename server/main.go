package main

import (
	"fmt"
	"insight/component"
	"insight/restful"
	"log"
	_ "net/http/pprof"
	"sync"

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
		webServer, err := restful.NewRestfulServer()
		if err != nil {
			fmt.Println()
		}
		webServer.Startup()
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		srv := component.NewHubServer()
		log.Fatal(gnet.Run(srv, srv.Addr, gnet.WithMulticore(srv.Multicore)))
	}()

	wg.Wait()

}
