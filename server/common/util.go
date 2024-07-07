package common

import "encoding/json"

func ToString(entry interface{}) string {
	bytes, _ := json.Marshal(entry)
	return string(bytes)
}
