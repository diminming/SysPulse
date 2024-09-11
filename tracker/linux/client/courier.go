package client

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"time"

	mutual_common "github.com/syspulse/mutual/common"
	"github.com/syspulse/tracker/linux/common"
)

type OnConnected func()

type Courier struct {
	SrvAddr string
	SrvPort int
	conn    net.Conn
}

func NewCourier() *Courier {
	submission := new(Courier)
	submission.SrvAddr = common.SysArgs.Server.Hub.Host
	submission.SrvPort = common.SysArgs.Server.Hub.Port

	for {
		flag := submission.Connect()
		if flag {
			break
		}
		time.Sleep(10 * time.Second)
	}

	return submission
}

func (c *Courier) Connect() bool {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", c.SrvAddr, c.SrvPort))
	// conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", c.SrvAddr, c.SrvPort), 30*time.Second)
	if err != nil {
		log.Default().Println("Error connecting:", err)
		return false
	}
	conn.(*net.TCPConn).SetWriteBuffer(1024 * 1024 * 10)
	c.conn = conn
	return true

}

func (c *Courier) Close() {
	c.conn.Close()
}

func (c *Courier) Write(payload []byte) {
	length := uint32(len(payload))
	log.Default().Printf("length of payload: %d", length)

	buffer := bytes.NewBuffer([]byte{})
	buffer.WriteByte('S')
	binary.Write(buffer, binary.LittleEndian, length)
	buffer.Write(payload)

	data := buffer.Bytes()
	md5 := mutual_common.MD5Calc(payload)
	log.Default().Printf("the md5 of payload_1: %s", md5)
	_, err := c.conn.Write(data)
	log.Default().Printf("payload: %s - Done.", md5)
	if err != nil {
		log.Default().Println(err)
	}
}
