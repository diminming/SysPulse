package network

import (
	"fmt"
	"log"
	"os"
	"path"
	"strings"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"github.com/google/gopacket/pcapgo"
	"github.com/google/uuid"

	"syspulse/tracker/linux/common"
	"syspulse/tracker/linux/task"
)

func GetPacket(ifName string, limit int64, filterSetting map[string]interface{}, successNotification func(), onFinish func(fPath string)) {
	tempDir := common.SysArgs.Storage.TempDir
	fPath := path.Join(tempDir, uuid.NewString())
	f, _ := os.Create(fPath)
	// 创建一个writer对象
	w := pcapgo.NewWriter(f)
	// 写入文件头，必须在调用前调用
	w.WriteFileHeader(uint32(65536), layers.LinkTypeEthernet)
	defer f.Close()

	handle, err := pcap.OpenLive(ifName, 65536, true, pcap.BlockForever)
	if err != nil {
		log.Fatal(err)
	}
	defer handle.Close()
	cdtLst := make([]string, 0)
	cdtLst = append(cdtLst, "tcp")
	for key, val := range filterSetting {
		if key == "port" {
			cdtLst = append(cdtLst, fmt.Sprintf("port %d", val))
		}

		if key == "direction" {
			switch val {
			case "in":
				cdtLst = append(cdtLst, fmt.Sprintf("src %s", filterSetting["ip"]))
			case "out":
				cdtLst = append(cdtLst, fmt.Sprintf("dst %s", filterSetting["ip"]))
			case "in,out":
				cdtLst = append(cdtLst, fmt.Sprintf("host %s", filterSetting["ip"]))
			}
		}
	}

	err = handle.SetBPFFilter(strings.Join(cdtLst, " and "))
	if err != nil {
		log.Fatal(err)
	}

	successNotification()

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	packets := int64(1)
	for packet := range packetSource.Packets() {
		w.WritePacket(packet.Metadata().CaptureInfo, packet.Data())
		packets = packets + 1
		if packets > limit {
			break
		}
	}
	onFinish(fPath)
}

func CreateCollectingTask(job task.Job) {

	ifName := job.IfName
	limit := job.Count
	direction := job.Direction

	go GetPacket(ifName, limit, map[string]interface{}{
		"port":      job.Port,
		"direction": strings.Join(direction, ","),
		"ip":        job.IpAddr,
	}, func() {
		task.UpdateJobStatus(job.Id, task.JOB_STATUS_RUNNING)
	}, func(fPath string) {
		bkName, objName := task.UploadOutcome(fPath)
		task.SendResult(job.Id, map[string]interface{}{
			"bucket": bkName,
			"object": objName,
		})
	})

}
