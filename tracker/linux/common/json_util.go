package common

import (
	"encoding/json"
)

func Stringify(data interface{}) string {
	b, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	s := string(b)
	return s
}
