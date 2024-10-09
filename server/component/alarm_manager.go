package component

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/syspulse/common"
	"github.com/syspulse/model"

	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/vm"
	"gopkg.in/yaml.v2"
)

const (
	CACHE_ITEM_LIMIT int64 = 200
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

func getData(key, field string, length int64) []float64 {
	listKey := key + ":" + field
	lst := model.CacheLRange(listKey, 0, length)
	result := make([]float64, 0, len(lst))
	for _, item := range lst {
		val, err := strconv.ParseFloat(item, 64)
		if err != nil {
			log.Default().Println("error convert item in value cache to float64: ", err)
			continue
		}
		result = append(result, val)
	}

	return result
}

func getFunctions() []expr.Option {
	return []expr.Option{
		expr.Env(model.PerfData{}),
		expr.Function(
			"endmost",
			func(params ...any) (any, error) {
				size := params[0].(int)
				target := params[1].(string)
				item := params[2].(string)

				return getData(target, item, int64(size)), nil
			},
		),
		expr.Function(
			"average",
			func(params ...any) (any, error) {
				array := params[0].([]any)
				total := float64(0)
				length := float64(len(array))
				for _, item := range array {
					switch val := item.(type) {
					case int:
						item1 := float64(val)
						total += item1
					case int32:
						item1 := float64(val)
						total += item1
					case int64:
						item1 := float64(val)
						total += item1
					case float32:
						item1 := float64(val)
						total += item1
					case float64:
						total += val
					}
				}
				return total / length, nil
			},
		),
		expr.Function(
			"maxItem",
			func(params ...any) (any, error) {
				array := params[0].([]any)
				max := float64(-1)
				for _, item := range array {
					switch val := item.(type) {
					case int:
						item1 := float64(val)
						if item1 > max {
							max = item1
						}
					case int32:
						item1 := float64(val)
						if item1 > max {
							max = item1
						}
					case int64:
						item1 := float64(val)
						if item1 > max {
							max = item1
						}
					case float32:
						item1 := float64(val)
						if item1 > max {
							max = item1
						}
					case float64:
						if val > max {
							max = val
						}
					}
				}
				return max, nil
			},
		), expr.Function(
			"minItem",
			func(params ...any) (any, error) {
				array := params[0].([]any)
				min := float64(-1)
				for _, item := range array {
					switch val := item.(type) {
					case int:
						item1 := float64(val)
						if item1 < min {
							min = item1
						}
					case int32:
						item1 := float64(val)
						if item1 < min {
							min = item1
						}
					case int64:
						item1 := float64(val)
						if item1 < min {
							min = item1
						}
					case float32:
						item1 := float64(val)
						if item1 < min {
							min = item1
						}
					case float64:
						item1 := float64(val)
						if item1 < min {
							min = item1
						}
					}
				}
				return min, nil
			},
		),
		expr.AsBool(),
	}
}

func CreateCacheItem(trigger *Trigger) (string, map[string]any) {
	id := trigger.Id
	express := trigger.Expre
	expObj, err := expr.Compile(express, getFunctions()...)

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
	// 初始化默认全局级别trigger设置
	defaultTriggerSetting := map[string]map[string]any{}
	for _, trigger := range TriggerConfig.Default {
		id, cacheItem := CreateCacheItem(trigger)
		defaultTriggerSetting[id] = cacheItem
	}

	// 初始化主机级别trigger设置
	for _, item := range TriggerConfig.TriggerSettingLst {

		identity := item.Identity
		triggerLst := item.TriggerLst
		// copy from global trigger setting
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

func CreateAlarmRecord(timestamp int64, linux *model.Linux, trigger_id, trigger, msg string) {

	sql := "insert into alarm(`timestamp`, `time_tag`, `linux_id`, `biz_id`, `trigger_id`, `trigger`,`ack`,`msg`,`create_timestamp`) value(?,?,?,?,?,?,?,?,?)"
	timeTag := timestamp2timeTag(timestamp)
	model.DBInsert(sql, timestamp, timeTag, linux.Id, linux.Biz.Id, trigger_id, trigger, false, msg, time.Now().UnixMilli())
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

func cleanCache(key string) {
	length := model.CacheLLen(key)
	if length > CACHE_ITEM_LIMIT {
		times := length - CACHE_ITEM_LIMIT
		for i := int64(0); i < times; i++ {
			model.CacheRPop(key)
		}
	}
}

func PutValue2Cache(data any, parentField string, identity string) {
	t := reflect.TypeOf(data)
	v := reflect.ValueOf(data)

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)

		if !value.IsValid() {
			continue
		}

		var fieldName string
		if parentField == "" {
			fieldName = field.Tag.Get("expr")
		} else {
			fieldName = parentField + "." + field.Tag.Get("expr")
		}
		if value.Kind() == reflect.Struct {
			PutValue2Cache(value.Interface(), fieldName, identity)
		} else {
			if !reflect.DeepEqual(value.Interface(), reflect.Zero(value.Type()).Interface()) {
				key := identity + ":" + fieldName
				// lockKey := "LCK_" + key

				// id, err := uuid.NewUUID()
				// if err != nil {
				// 	log.Default().Println("error lock: ", err)
				// }
				// model.AcquireLock(lockKey, id.String(), time.Second*1)
				// defer model.ReleaseLock(lockKey, id.String())

				length := model.CacheLPUSH(key, value.Interface())
				// 如果对象长度超过阈值则进行清理，减少内存消耗
				if length > CACHE_ITEM_LIMIT {
					cleanCache(key)
				}
			}
		}
	}
}

func TriggerCheck(identity string, data model.PerfData, dataType model.PerformenceDataType, timestamp int64) {
	triggerSettings := TriggerCache[identity]

	// 使用反射将对象的值放入缓存
	PutValue2Cache(data, "", identity)

	for triggerId, trigger := range triggerSettings {
		express := trigger["exp"].(string)

		program := trigger["expObj"].(*vm.Program)
		result, err := expr.Run(program, data)

		if err != nil {
			log.Default().Panicf("error calc at target: %s, exp: %s, data: %s. \n%v", identity, program.Source().String(), common.ToString(data), err)
		}

		if result.(bool) && !CheckTargetIsActive(identity, triggerId) {
			linux := model.LoadLinuxByIdentity(identity)

			msg := CreateMessage(trigger, linux)
			CreateAlarmRecord(timestamp, linux, triggerId, express, msg)
		}

	}
}
