package model

type Menu struct {
	ID              int64         `json:"id"`
	Title           string        `json:"title"`
	Identity        string        `json:"identity"`
	Type            string        `json:"type"`
	Index           uint8         `json:"index"`
	ParentId        int64         `json:"parentId"`
	ParentTitle     string        `json:"parentTitle"`
	Url             string        `json:"url"`
	PermissionLst   []*Permission `json:"permissionLst"`
	CreateTimestamp int64         `json:"createTimestamp"`
	UpdateTimestamp int64         `json:"updateTimestamp"`
}

type Permission struct {
	ID              int64  `json:"id"`
	Identity        string `json:"identity"`
	Name            string `json:"name"`
	Url             string `json:"url"`
	Method          string `json:"method"`
	CreateTimestamp int64  `json:"createTimestamp"`
	UpdateTimestamp int64  `json:"updateTimestamp"`
}

type Role struct {
	ID              int64         `json:"id"`
	Name            string        `json:"name"`
	Identity        string        `json:"identity"`
	PermissionLst   []*Permission `json:"permissionLst"`
	CreateTimestamp int64         `json:"createTimestamp"`
	UpdateTimestamp int64         `json:"updateTimestamp"`
}

type User struct {
	ID              int64   `json:"id"`
	Username        string  `json:"username"`
	Passwd          string  `json:"password"`
	IsActive        bool    `json:"isActive"`
	RoleLst         []*Role `json:"roleLst"`
	CreateTimestamp int64   `json:"createTimestamp"`
	UpdateTimestamp int64   `json:"updateTimestamp"`
}

func GetTotalOfPermission() int64 {
	total := int64(0)
	sqlstr := "select count(id) as total from permission"
	result := SqlDB.QueryRow(sqlstr)
	result.Scan(&total)
	return total
}
