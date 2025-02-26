package handler

import (
	"bufio"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/syspulse/common"
	"github.com/syspulse/component"
	"github.com/syspulse/model"
	"github.com/syspulse/restful/server/response"
	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

func sendJob2Agent(job *model.Job) error {
	linuxId := job.LinuxId
	linux := GetLinuxById(linuxId)

	reqBody, err := json.Marshal(job)
	if err != nil {
		log.Default().Print("Error converting job info: ", err)
		return err
	}
	resp, err := http.Post(fmt.Sprintf("http://%s/api/job", linux.AgentConn), "application/json", strings.NewReader(string(reqBody)))
	if err != nil {
		log.Default().Print("Error sending request: ", err)
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Default().Print("Error reading response: ", err)
		return err
	}

	if resp.StatusCode != 200 {
		panic(fmt.Sprintf("Error creating job: %d, %s", resp.StatusCode, string(body)))
	}
	return nil
}

func CreateAnalJob(job *model.Job) (*model.Job, error) {

	if job.LinuxId == 0 {
		panic("error input: linux id is 0")
	}

	job.Status = model.JOB_STATUS_CREATED
	job.CreateTimestamp = time.Now().UnixMilli()
	job.UpdateTimestamp = time.Now().UnixMilli()

	tx, err := model.SqlDB.Begin()
	if err != nil {
		log.Default().Println(err)
	}
	defer tx.Rollback()
	sql := "insert into job(`job_name`, `category`, `type`, `status`, `startup_time`, `linux_id`, `pid`, `duration`, `immediately`, `create_timestamp`, `update_timestamp`) value(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
	result, err := tx.Exec(sql, job.JobName, job.Category, job.Type, job.Status, job.StartupTime, job.LinuxId, job.Pid, job.Duration, job.Immediately, job.CreateTimestamp, job.UpdateTimestamp)
	if err != nil {
		log.Default().Println(err)
	}
	insertId, err := result.LastInsertId()
	if err != nil {
		log.Default().Println(err)
	}
	job.Id = insertId

	err = sendJob2Agent(job)

	if err != nil {
		tx.Rollback()
		panic(fmt.Errorf("error create job on agent: %v", err))
	}

	err = tx.Commit()
	if err != nil {
		panic(fmt.Errorf("error commit job info 2 DB: %v", err))
	}

	return job, nil

}

func CreateTrafficJob(job *model.Job) (*model.Job, error) {
	log.Default().Println(common.Stringfy(job))
	if job.LinuxId == 0 {
		panic("error input: linux id is 0")
	}

	job.Status = model.JOB_STATUS_CREATED
	job.CreateTimestamp = time.Now().UnixMilli()
	job.UpdateTimestamp = time.Now().UnixMilli()

	tx, err := model.SqlDB.Begin()
	if err != nil {
		log.Default().Println(err)
	}
	defer tx.Rollback()

	sql := "insert into job(`job_name`, `category`, `status`, `direction`, `count`,`ip_addr`, `if_name`,`port`,`linux_id`,`create_timestamp`,`update_timestamp`) value(?,?,?,?,?,?,?,?,?,?,?)"
	result, err := tx.Exec(sql, job.JobName, job.Category, job.Status, strings.Join(job.Direction, ","), job.Count, job.IpAddr, job.IfName, job.Port, job.LinuxId, job.CreateTimestamp, job.UpdateTimestamp)
	if err != nil {
		log.Default().Println(err)
	}
	insertId, err := result.LastInsertId()
	if err != nil {
		log.Default().Println(err)
	}
	job.Id = insertId

	err = sendJob2Agent(job)

	if err != nil {
		tx.Rollback()
		panic(fmt.Errorf("error create job on agent: %v", err))
	}

	err = tx.Commit()
	if err != nil {
		panic(fmt.Errorf("error commit job info 2 DB: %v", err))
	}

	return job, nil
}

func Write2File(filename string, data []byte) {
	file, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// 使用缓冲区写入数据
	writer := bufio.NewWriter(file)
	_, err = writer.Write(data)
	if err != nil {
		panic(err)
	}
	writer.Flush()

	println("Binary data has been written to", filename)
}

func ReadFromFile(filename string) (string, error) {
	// 打开文件
	file, err := os.Open(filename)
	if err != nil {
		log.Default().Println(err)
		return "", err
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	data := make([]byte, 0, 64*1024)
	// 创建一个字节切片用于存储数据
	buffer := make([]byte, 1024) // 选择合适的缓冲区大小

	// 从文件中读取数据
	for {
		// 从 Reader 中读取数据到缓冲区
		n, err := reader.Read(buffer)
		if err != nil {
			break
		}
		// 处理读取的数据
		data = append(data, buffer[:n]...)
	}

	return string(data), nil

}

func CreateJob(ctx *gin.Context) {
	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		log.Default().Println(err)
		return
	}
	job := new(model.Job)
	err = json.Unmarshal(body, job)
	if err != nil {
		log.Default().Println(err)
		panic(err)
	}

	if job.Category == "proc_profiling" {
		CreateAnalJob(job)
	} else if job.Category == "traffic" {
		CreateTrafficJob(job)
	}

	ctx.JSON(http.StatusOK, response.JsonResponse{Status: http.StatusOK, Data: job, Msg: "success"})
}

func UpdateJobStatus(ctx *gin.Context) {
	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		log.Default().Println(err)
		return
	}
	data := new(map[string]interface{})
	err = json.Unmarshal(body, data)
	if err != nil {
		log.Default().Println(err)
		return
	}
	jobId := int64((*data)["jobId"].(float64))
	status := int32((*data)["status"].(float64))
	model.DBUpdate("update job set `status` = ? where `id` = ?", status, jobId)
	ctx.JSON(http.StatusOK, response.JsonResponse{Status: http.StatusOK, Msg: "ok"})
}

func GetTrafficJobCountByLinuxId(linuxId int64) int64 {
	sql := "select count(id) as count from `job` where `linux_id`=? and category='traffic'"
	lst := model.DBSelect(sql, linuxId)
	for _, o := range lst {
		return o["count"].(int64)
	}
	return 0
}

func GetTrafficJobLstByLinuxId(linuxId int64, page int64, pageSize int64) []model.Job {
	first := page * pageSize
	sql := "select `id`, `job_name`, `status`, `direction`, `count`, `port`, `if_name`, `ip_addr`, `create_timestamp`, `update_timestamp` from `job` where `linux_id`=? and category='traffic' order by `update_timestamp` desc limit ?, ?"
	lst := model.DBSelect(sql, linuxId, first, pageSize)
	result := []model.Job{}

	for _, o := range lst {

		item := model.Job{
			Id:              o["id"].(int64),
			JobName:         string(o["job_name"].([]uint8)),
			IfName:          string(o["if_name"].([]uint8)),
			IpAddr:          string(o["ip_addr"].([]uint8)),
			Direction:       strings.Split(string(o["direction"].([]uint8)), ","),
			Count:           o["count"].(int64),
			Port:            int32(o["port"].(int64)),
			Status:          int(o["status"].(int64)),
			CreateTimestamp: o["create_timestamp"].(int64),
			UpdateTimestamp: o["update_timestamp"].(int64),
		}
		result = append(result, item)
	}
	return result
}

func GetJobById(jobId int64) *model.Job {
	sql := "select id, job_name, category,extend, create_timestamp, update_timestamp from job where id = ?"
	lst := model.DBSelect(sql, jobId)
	job := new(model.Job)

	for _, item := range lst {
		job.Id = item["id"].(int64)
		job.JobName = string(item["job_name"].([]uint8))
		job.Category = string(item["category"].([]uint8))
		job.CreateTimestamp = item["create_timestamp"].(int64)
		job.UpdateTimestamp = item["update_timestamp"].(int64)
		if item["extend"] != nil {
			job.Extend = string(item["extend"].([]uint8))
		}
		return job
	}

	return nil
}

func OnProfilingJobFinished(jobId int64, body []byte) {
	data := new([]map[string]interface{})
	err := json.Unmarshal(body, data)
	if err != nil {
		log.Default().Println(err)
		return
	}
	filename := fmt.Sprintf("%s/insight_%d.json", common.SysArgs.Storage.File.Path, jobId)

	Write2File(filename, body)

	model.DBUpdate("update job set `status` = ? where `id` = ?", model.JOB_STATUS_FINISHED, jobId)
}

func OnTrafficJobFinished(jobId int64, body []byte) {
	data := new(map[string]interface{})
	err := json.Unmarshal(body, data)
	if err != nil {
		log.Default().Println(err)
		return
	}
	result := *data
	extend := fmt.Sprintf("%s:%s", result["bucket"], result["object"])
	log.Default().Println(data, extend)

	model.DBUpdate("update job set `status`=?, `extend`=? where `id`=?", model.JOB_STATUS_FINISHED, extend, jobId)
}

func GetJobCount(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, response.JsonResponse{Status: http.StatusOK, Msg: "success", Data: model.GetJobTotal()})
}

func OnJobFinished(ctx *gin.Context) {

	jobId, err := strconv.ParseInt(ctx.Param("jobId"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.JsonResponse{Status: http.StatusBadRequest, Msg: "job id is not a number."})
	}

	job := GetJobById(jobId)

	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		log.Default().Println(err)
		return
	}

	switch job.Category {
	case "proc_profiling":
		OnProfilingJobFinished(jobId, body)
	case "traffic":
		OnTrafficJobFinished(jobId, body)
	}

	ctx.JSON(http.StatusOK, response.JsonResponse{Status: http.StatusOK, Msg: "ok"})
}

func GetTrafficJobLst(ctx *gin.Context) {
	values := ctx.Request.URL.Query()

	page, err := strconv.ParseInt(values.Get("page"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.JsonResponse{Status: http.StatusBadRequest, Msg: "page is not a number."})
		return
	}

	pageSize, err := strconv.ParseInt(values.Get("pageSize"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.JsonResponse{Status: http.StatusBadRequest, Msg: "pageSize is not a number."})
		return
	}

	linuxId, err := strconv.ParseInt(values.Get("linuxId"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.JsonResponse{Status: http.StatusBadRequest, Msg: "linuxId is not a number."})
		return
	}
	result := GetTrafficJobLstByLinuxId(linuxId, page, int64(pageSize))
	total := GetTrafficJobCountByLinuxId(linuxId)

	ctx.JSON(http.StatusOK, response.JsonResponse{Status: http.StatusOK, Msg: "ok", Data: map[string]interface{}{
		"lst":   result,
		"total": total,
	}})
}

func getTCPFlags(tcp *layers.TCP) []string {
	flags := make([]string, 0)
	if tcp.SYN {
		flags = append(flags, "SYN")
	}
	if tcp.ACK {
		flags = append(flags, "ACK")
	}
	if tcp.FIN {
		flags = append(flags, "FIN")
	}
	if tcp.PSH {
		flags = append(flags, "PSH")
	}
	if tcp.URG {
		flags = append(flags, "URG")
	}
	if tcp.RST {
		flags = append(flags, "RST")
	}
	return flags
}

type Connection struct {
	PacketLst []*layers.TCP
}

func (conn *Connection) GetPreviousPacket() *layers.TCP {
	return conn.PacketLst[len(conn.PacketLst)-1]
}

func (conn *Connection) PushPacket(packet *layers.TCP) {
	conn.PacketLst = append(conn.PacketLst, packet)
}

func ParsePcapFile(filePath string) ([]*Connection, uint32) {

	f, err := os.Open(filePath)
	if err != nil {
		zap.L().Panic("can't open pcap file.", zap.String("file", filePath))
	}
	defer f.Close()

	handle, err := pcap.OpenOfflineFile(f)
	if err != nil {
		zap.L().Panic("error open pcap file.", zap.Error(err))
	}
	defer handle.Close()

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	connLst := make([]*Connection, 0)
	count := uint32(0)
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
					count += 1
					tcp, _ := transportLayer.(*layers.TCP)
					// 提取TCP端口信息
					srcPort := tcp.SrcPort.String()
					dstPort := tcp.DstPort.String()
					flags := getTCPFlags(tcp)
					log.Default().Printf(">>> %s:%s -> %s:%s; len: %d; flags: %s; seq: %d; ack: %d\n", src, srcPort, dst, dstPort, length, strings.Join(flags, ","), tcp.Seq, tcp.Ack)

					if tcp.SYN && !tcp.ACK {
						conn := new(Connection)
						conn.PushPacket(tcp)
						connLst = append(connLst, conn)
					} else if tcp.SYN && tcp.ACK {
						for _, conn := range connLst {
							prev := conn.GetPreviousPacket()
							if prev.SYN && (prev.Seq+1 == tcp.Ack) {
								conn.PushPacket(tcp)
							}
						}
					} else if tcp.PSH {
						for _, conn := range connLst {
							prev := conn.GetPreviousPacket()
							if prev.ACK && (tcp.Seq == prev.Seq) {
								conn.PushPacket(tcp)
							}
						}
					} else if tcp.FIN {
						for _, conn := range connLst {
							prev := conn.GetPreviousPacket()
							if prev.ACK && (tcp.Ack == prev.Seq) {
								conn.PushPacket(tcp)
							} else if prev.ACK && (tcp.Ack == prev.Seq+1) {
								conn.PushPacket(tcp)
							}
						}
					} else if tcp.ACK {
						for _, conn := range connLst {
							prev := conn.GetPreviousPacket()
							if prev.SYN && (tcp.Ack == prev.Seq+1) {
								conn.PushPacket(tcp)
							} else if prev.PSH && (tcp.Ack == prev.Seq+uint32(len(prev.Payload))) {
								conn.PushPacket(tcp)
							} else if prev.FIN && (tcp.Ack == prev.Seq+1) {
								conn.PushPacket(tcp)
							} else if prev.ACK && (tcp.Ack == prev.Seq) {
								conn.PushPacket(tcp)
							}
						}
					}

				}
			}
		}
	}
	return connLst, count
}

var PATTERN_URL_IN_PAYLOAD = regexp.MustCompile(`(?m)^(GET|PATCH|POST|PUT|DELETE|HEAD|CONNECT|OPTIONS|TRACE) ([\w\/\?\=\&]+) (HTTP\/[\d\.]+)\r$`)
var PATTERN_STATUS_IN_PAYLOAD = regexp.MustCompile(`(?m)^(HTTP/[\d\.]+) \d{3} \w+\r$`)

func GetTimestamp(packet *layers.TCP) (uint32, uint32) {
	opts := packet.Options
	for _, opt := range opts {
		if opt.OptionType == layers.TCPOptionKindTimestamps {
			timestamp := binary.BigEndian.Uint32(opt.OptionData[0:4])
			timestampEcho := binary.BigEndian.Uint32(opt.OptionData[4:8])
			return timestamp, timestampEcho
		}
	}
	return uint32(0), uint32(0)
}

func Transform(connLst []*Connection) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, 10)
	for _, conn := range connLst {
		item := make(map[string]any)
		packetLst := make([]map[string]any, 0)
		throughput := uint64(0)
		for _, packet := range conn.PacketLst {
			p := make(map[string]any)
			p["Seq"] = packet.Seq
			p["Ack"] = packet.Ack
			p["SrcPrt"] = packet.SrcPort
			p["DstPrt"] = packet.DstPort

			timestamp, timestampEcho := GetTimestamp(packet)
			p["timestamp"] = timestamp
			p["timestampEcho"] = timestampEcho

			p["flags"] = strings.Join(getTCPFlags(packet), ", ")

			packetLst = append(packetLst, p)

			if packet.PSH {
				throughput += uint64(len(packet.Payload))
				payload := string(packet.Payload)
				req1 := PATTERN_URL_IN_PAYLOAD.FindString(payload)
				status1 := PATTERN_STATUS_IN_PAYLOAD.FindString(payload)
				if req1 != "" {
					item["req"] = strings.TrimSpace(req1)
				}
				if status1 != "" {
					item["status"] = strings.TrimSpace(status1)
					item["time"] = timestamp - packetLst[0]["timestamp"].(uint32)
				}
			}
		}
		item["throughput"] = throughput
		item["packet_lst"] = packetLst
		result = append(result, item)
	}
	return result
}

func TrafficJobHandler(job *model.Job) map[string]interface{} {
	info := strings.Split(job.Extend, ":")
	if len(info) != 2 {
		log.Default().Panicln("job.Extend is wrong: ", job.Extend)
	}
	filePath := component.DownloadFromFileServer(info[0], info[1])
	res, count := ParsePcapFile(filePath)
	return map[string]any{
		"lst":          Transform(res),
		"count_packet": count,
		"count_conn":   len(res),
	}
}

func GetJobResult(ctx *gin.Context) {
	jobId, err := strconv.ParseInt(ctx.Param("jobId"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.JsonResponse{Status: http.StatusBadRequest, Msg: "job id is not a number."})
		return
	}
	job := GetJobById(jobId)
	switch job.Category {
	case "proc_profiling":
		filename := fmt.Sprintf("%s/insight_%d.json", common.SysArgs.Storage.File.Path, jobId)
		ctn, err := ReadFromFile(filename)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, response.JsonResponse{Status: http.StatusInternalServerError, Msg: err.Error()})
			return
		}
		// 解析JSON字符串到map
		var data []map[string]interface{}
		err = json.Unmarshal([]byte(ctn), &data)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, response.JsonResponse{Status: http.StatusInternalServerError, Msg: err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, response.JsonResponse{Status: http.StatusOK, Msg: "ok", Data: data})
	case "traffic":
		data := TrafficJobHandler(job)
		ctx.JSON(http.StatusOK, response.JsonResponse{Status: http.StatusOK, Msg: "ok", Data: data})
	}

}

func DeleteJob(ctx *gin.Context) {
	jobId, err := strconv.ParseInt(ctx.Param("jobId"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.JsonResponse{Status: http.StatusBadRequest, Msg: "job id is not a number."})
		return
	}
	count := model.DeleteJob(jobId)
	ctx.JSON(http.StatusOK, response.JsonResponse{Status: http.StatusOK, Msg: "ok", Data: count})
}
