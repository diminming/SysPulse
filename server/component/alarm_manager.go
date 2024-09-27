package component

import (
	"log"
	"os"
	"time"

	"github.com/syspulse/common"
	"github.com/syspulse/model"

	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/vm"
	"gopkg.in/yaml.v2"
)

type Trigger struct {
	LinuxIdentity string   `yaml:"linux"`
	Expres        []string `yaml:"expression"`
}

type TriggerSetting struct {
	TriggerLst []*Trigger `yaml:"triggers"`
}

var TriggerConfig TriggerSetting

var TriggerCache map[string][]*vm.Program

func init() {

	yamlFile, err := os.ReadFile(common.SysArgs.TriggerCfg)
	if err != nil {
		log.Default().Fatalf("can't open config file: %v", err)
	}
	err = yaml.Unmarshal(yamlFile, &TriggerConfig)
	if err != nil {
		log.Fatalf("can't read config file: %v", err)
	}
	TriggerCache = make(map[string][]*vm.Program)
	buildTrigger()
}

func buildTrigger() {

	for _, item := range TriggerConfig.TriggerLst {
		identity := item.LinuxIdentity
		expressions := item.Expres
		expObjLst := make([]*vm.Program, 0, len(expressions))
		for _, expression := range expressions {
			expObj, err := expr.Compile(expression, expr.Env(model.PerfData{}), expr.AsBool())
			if err != nil {
				log.Default().Fatalf("can't create expression by %s.\n%v\n", expression, err)
			}
			expObjLst = append(expObjLst, expObj)
		}
		TriggerCache[identity] = expObjLst
	}
}

func timestamp2timeTag(timestmap int64) string {
	return time.UnixMilli(timestmap).Format("2006010215")
}

func createAlarmRecord(timestamp int64, identity string, trigger string, parameters model.PerfData) {
	linux := model.LoadLinuxByIdentity(identity)
	// linuxId := model.CacheGet(identity)
	perfDataStr := common.ToString(parameters)

	sql := "insert into alarm(`timestamp`, `time_tag`, `linux_id`, `biz_id`, `trigger`,`ack`,`perf_data`,`create_timestamp`) value(?,?,?,?,?,?,?,?)"
	timeTag := timestamp2timeTag(timestamp)
	model.DBInsert(sql, timestamp, timeTag, linux.Id, linux.Biz.Id, trigger, false, []byte(perfDataStr), time.Now().UnixMilli())
}

func TriggerCheck(identity string, parameters model.PerfData, timestamp int64) {
	programs := TriggerCache[identity]
	for _, program := range programs {

		result, err := expr.Run(program, parameters)

		if err != nil {
			log.Default().Panicf("error calc at linux: %s, exp: %s, data: %s. \n%v", identity, program.Source().String(), common.ToString(parameters), err)
		}

		if result.(bool) {
			log.Default().Printf("<Alarm> timestamp: %d, identity: %s, result: %b, trigger: %s", timestamp, identity, result, program.Source().String())
			createAlarmRecord(timestamp, identity, program.Source().String(), parameters)
		}

	}
}
