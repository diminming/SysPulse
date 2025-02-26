package common

import "encoding/json"

func Stringfy(entry interface{}) string {
	bytes, _ := json.Marshal(entry)
	return string(bytes)
}
