package client

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"time"

	"github.com/syspulse/tracker/linux/common"
	"go.uber.org/zap"
)

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
		zap.L().Error("Error connecting:", zap.Error(err))
		return false
	}
	conn.(*net.TCPConn).SetWriteBuffer(1024 * 1024 * 10)
	c.conn = conn
	return true

}

func (c *Courier) Close() {
	c.conn.Close()
}

func Pack(payload []byte) []byte {
	buffer := bytes.NewBuffer([]byte{})
	buffer.WriteByte('S')

	length := uint32(len(payload))
	binary.Write(buffer, binary.LittleEndian, length)

	buffer.Write(payload)

	return buffer.Bytes()

}

func (c *Courier) Write(payload []byte) {
	data := Pack(payload)

	// md5 := mutual_common.MD5Calc(payload)
	// log.Default().Printf("the md5 of payload: %s", md5)
	for {
		_, err := c.conn.Write(data)
		// log.Default().Printf("payload: %s - Done.", md5)
		if err != nil {
			zap.L().Error("error write data to connection: ", zap.Error(err))
			reConnected := c.Connect()
			if reConnected {
				break
			}
			time.Sleep(3 * time.Second)
		} else {
			break
		}
	}

}
