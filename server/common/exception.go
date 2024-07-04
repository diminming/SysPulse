package common

import "fmt"

type InsightException struct {
	Code int
	Msg  string
}

func (e *InsightException) Error() string {
	return fmt.Sprintf("Error %d: %s", e.Code, e.Msg)
}
