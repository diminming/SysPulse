package client

import (
	"fmt"
	"log"
	"net"
	"syspulse/tracker/linux/common"
	"time"
)

type OnConnected func()

type Courier struct {
	SrvAddr   string
	SrvPort   int
	send_pipe chan []byte
	conn      net.Conn
}

func NewCourier() *Courier {
	submission := new(Courier)
	submission.SrvAddr = common.SysArgs.Server.Hub.Host
	submission.SrvPort = common.SysArgs.Server.Hub.Port
	submission.send_pipe = make(chan []byte, 10)

	callback := func() {
		go submission.Write()
	}

	for {
		flag := submission.Connect(callback)
		if flag {
			break
		}
		time.Sleep(10 * time.Second)
	}

	return submission
}

func (c *Courier) Connect(callback OnConnected) bool {
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", c.SrvAddr, c.SrvPort), 30*time.Second)
	if err != nil {
		log.Default().Println("Error connecting:", err)
		return false
	}
	c.conn = conn
	if callback != nil {
		callback()
	}
	return true

}

func (c *Courier) Send(data []byte) {
	c.send_pipe <- data
	// log.Default().Printf("length of c.send_pipe: %d", len(c.send_pipe))
}

func (c *Courier) Close() {
	c.conn.Close()
}

func (c *Courier) Write() {
	for {
		data := <-c.send_pipe
		_, err := c.conn.Write(data)
		if err != nil {
			log.Default().Println(err)
			c.Connect(nil)
		}
	}
}
