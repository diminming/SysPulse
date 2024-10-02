package component

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/syspulse/common"
	"github.com/syspulse/model"

	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/vm"
	"gopkg.in/yaml.v2"
)

type Trigger struct {
	Id    string `yaml:"id"`
	Msg   string `yaml:"message"`
	Expre string `yaml:"expression"`
}

type TriggerSetting struct {
	Identity   string     `yaml:"identity"`
	TriggerLst []*Trigger `yaml:"triggers"`
}

var TriggerConfig struct {
	Default           []*Trigger        `yaml:"_default"`
	TriggerSettingLst []*TriggerSetting `yaml:"trigger_setting"`
}

var TriggerCache map[string]map[string]map[string]any

func init() {

	yamlFile, err := os.ReadFile(common.SysArgs.TriggerCfg)
	if err != nil {
		log.Default().Fatalf("can't open config file: %v", err)
	}
	err = yaml.Unmarshal(yamlFile, &TriggerConfig)
	if err != nil {
		log.Fatalf("can't read config file: %v", err)
	}

	TriggerCache = make(map[string]map[string]map[string]any)

	buildTrigger()

}

func CreateCacheItem(trigger *Trigger) (string, map[string]any) {
	id := trigger.Id
	express := trigger.Expre
	expObj, err := expr.Compile(express, expr.Env(model.PerfData{}), expr.AsBool())

	if err != nil {
		log.Default().Fatalf("can't create expression by %s.\n%v\n", express, err)
	}

	return id, map[string]any{
		"id":     id,
		"exp":    express,
		"expObj": expObj,
		"msg":    trigger.Msg,
	}

}

func copyDefaultTriggerSetting(defaultTriggerSetting map[string]map[string]any) map[string]map[string]any {
	copyMap := make(map[string]map[string]any)
	for key, value := range defaultTriggerSetting {
		copyMap[key] = value
	}
	return copyMap
}

func buildTrigger() {
	defaultTriggerSetting := map[string]map[string]any{}
	for _, trigger := range TriggerConfig.Default {
		id, cacheItem := CreateCacheItem(trigger)
		defaultTriggerSetting[id] = cacheItem
	}

	for _, item := range TriggerConfig.TriggerSettingLst {

		identity := item.Identity
		triggerLst := item.TriggerLst

		defaultCopy := copyDefaultTriggerSetting(defaultTriggerSetting)

		for _, trigger := range triggerLst {
			id, cacheItem := CreateCacheItem(trigger)
			defaultCopy[id] = cacheItem
		}

		TriggerCache[identity] = defaultCopy

	}
	log.Default().Println("init trigger finished.")
}

func timestamp2timeTag(timestmap int64) string {
	return time.UnixMilli(timestmap).Format("2006010215")
}

func CreateAlarmRecord(timestamp int64, linux *model.Linux, trigger_id, trigger, perfDataStr string, msg string) {

	sql := "insert into alarm(`timestamp`, `time_tag`, `linux_id`, `biz_id`, `trigger_id`, `trigger`,`ack`,`msg`,`perf_data`,`create_timestamp`) value(?,?,?,?,?,?,?,?,?,?)"
	timeTag := timestamp2timeTag(timestamp)
	model.DBInsert(sql, timestamp, timeTag, linux.Id, linux.Biz.Id, trigger_id, trigger, false, msg, []byte(perfDataStr), time.Now().UnixMilli())
	key := fmt.Sprintf("alarm_%s", linux.LinuxId)
	model.CacheHSet(key, trigger, "true")

}

func CheckTargetIsActive(identity, trigger string) bool {
	key := fmt.Sprintf("alarm_%s", identity)
	value := model.CacheHGet(key, trigger)
	return value == "true"
}

var BUILDIN_MACROS = map[string]func(*model.Linux) string{
	"<!hostname>": func(linux *model.Linux) string {
		return linux.Hostname
	},
	"<!hostID>": func(linux *model.Linux) string {
		return linux.LinuxId
	},
}

func CreateMessage(trigger map[string]any, linux *model.Linux) string {
	msg := trigger["msg"].(string)
	for macro, handler := range BUILDIN_MACROS {
		msg = strings.ReplaceAll(msg, macro, handler(linux))
	}
	return msg
}

func TriggerCheck(identity string, parameters model.PerfData, dataType model.PerformenceDataType, timestamp int64) {
	triggerSettings := TriggerCache[identity]
	for triggerId, trigger := range triggerSettings {
		log.Default().Println("trigger id: ", triggerId)
		express := trigger["exp"].(string)

		if (dataType == model.DataType_CpuPerformence && strings.HasPrefix(express, "cpu.")) ||
			(dataType == model.DataType_MemoryPerformence && strings.HasPrefix(express, "memory.")) ||
			(dataType == model.DataType_LoadPerformence && strings.HasPrefix(express, "load.")) ||
			(dataType == model.DataType_SwapPerformence && strings.HasPrefix(express, "swap.")) {

			program := trigger["expObj"].(*vm.Program)
			result, err := expr.Run(program, parameters)
			if err != nil {
				log.Default().Panicf("error calc at target: %s, exp: %s, data: %s. \n%v", identity, program.Source().String(), common.ToString(parameters), err)
			}

			if result.(bool) && !CheckTargetIsActive(identity, triggerId) {
				linux := model.LoadLinuxByIdentity(identity)
				msg := CreateMessage(trigger, linux)
				perfDataStr := ""

				if strings.HasPrefix(express, "cpu.") {
					perfDataStr = common.ToString(parameters.CPU)
				} else if strings.HasPrefix(express, "memory.") {
					perfDataStr = common.ToString(parameters.Memory)
				} else if strings.HasPrefix(express, "load.") {
					perfDataStr = common.ToString(parameters.Load)
				} else if strings.HasPrefix(express, "swap.") {
					perfDataStr = common.ToString(parameters.Swap)
				}

				CreateAlarmRecord(timestamp, linux, triggerId, express, perfDataStr, msg)
			}

		}

	}
}
