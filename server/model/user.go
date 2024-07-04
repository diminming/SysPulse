package model

type User struct {
	ID              int64  `json:"id"`
	Username        string `json:"username"`
	Passwd          string `json:"passwd"`
	IsActive        bool   `json:"is_active"`
	CreateTimestamp int64  `json:"create_timestamp"`
	UpdateTimestamp int64  `json:"update_timestamp"`
}

const (
	JOB_STATUS_CREATED  = 1
	JOB_STATUS_RUNNING  = 2
	JOB_STATUS_FINISHED = 3
)

type Job struct {
	Id              int64  `json:"id"`
	JobName         string `json:"job_name"`
	Category        string `json:"category"`
	Type            string `json:"type"`
	Status          int    `json:"status"`
	StartupTime     int64  `json:"startup_time"`
	LinuxId         int64  `json:"linux_id"`
	Pid             int32  `json:"pid"`
	Duration        int32  `json:"duration"`
	Immediately     bool   `json:"immediately"`
	CreateTimestamp int64  `json:"create_timestamp"`
	UpdateTimestamp int64  `json:"update_timestamp"`
}
