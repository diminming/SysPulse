package net

import (
	"fmt"
	"log"
	"syspulse/tracker/linux/client"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

type Collector struct {
	DeviceLst string
	Courier   *client.Courier
}

func (collector *Collector) Run() {

}

type Gatherer struct {
	Device  string
	Snaplen int32
}

func NewCollector(courier *client.Courier) *Collector {
	collector := new(Collector)
	collector.Courier = courier
	return collector
}

func GetPacket() {
	handle, err := pcap.OpenLive("wlp0s20f3", 1600, true, pcap.BlockForever)
	if err != nil {
		log.Fatal(err)
	}
	defer handle.Close()

	// 设置过滤条件，这里设置为捕获TCP流量
	err = handle.SetBPFFilter("tcp")
	if err != nil {
		log.Fatal(err)
	}

	// 开始抓包
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())

	// 遍历捕获到的数据包
	for packet := range packetSource.Packets() {
		// 解析数据包
		ipv4Layer := packet.Layer(layers.LayerTypeIPv4)
		if ipv4Layer != nil {
			ipv4 := ipv4Layer.(*layers.IPv4)
			src := ipv4.SrcIP.String()
			dst := ipv4.DstIP.String()
			length := ipv4.Length

			transportLayer := packet.TransportLayer()
			if transportLayer != nil {
				// 检查传输层类型
				switch transportLayer.(type) {
				case *layers.TCP:
					tcp, _ := transportLayer.(*layers.TCP)
					// 提取TCP端口信息
					srcPort := tcp.SrcPort.String()
					dstPort := tcp.DstPort.String()
					fmt.Printf("TCP { src: [%s:%s], dst: [%s:%s], len: %d }\n", src, srcPort, dst, dstPort, length)
				case *layers.UDP:
					udp, _ := transportLayer.(*layers.UDP)
					// 提取UDP端口信息
					srcPort := udp.SrcPort.String()
					dstPort := udp.DstPort.String()
					fmt.Printf("UDP { src: [%s:%s], dst: [%s:%s], len: %d }\n", src, srcPort, dst, dstPort, length)
				}
			}

		}
	}
}
