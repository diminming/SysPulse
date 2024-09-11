package model

type Alarm struct {
	Id              int64
	Timestamp       int64
	CreateTimestamp int64
	Trigger         string
}
