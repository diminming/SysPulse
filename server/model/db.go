package model

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/syspulse/common"

	_ "github.com/go-sql-driver/mysql"
)

var SqlDB *sql.DB

func init() {

	dbCnf := common.SysArgs.Storage.DB

	SqlDB, _ = sql.Open(
		"mysql",
		fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local", // 使用本地时间，即东八区，北京时间
			dbCnf.User,
			dbCnf.Password,
			dbCnf.Host,
			dbCnf.Port,
			dbCnf.Database))
	// set pool params
	SqlDB.SetMaxOpenConns(100)
	SqlDB.SetMaxIdleConns(100)
	SqlDB.SetConnMaxLifetime(time.Minute * 60) // mysql default conn timeout=8h, should < mysql_timeout
	err := SqlDB.Ping()
	if err != nil {
		log.Fatal("database init failed, err: ", err)
	}
	log.Println("mysql conn pool has initiated.")
}

func DBSelectRow(sql string, args ...interface{}) map[string]interface{} {
	lst := DBSelect(sql, args...)
	length := len(lst)
	if length == 1 {
		return lst[0]
	} else if length == 0 {
		log.Default().Println("no result...")
	} else {
		log.Default().Println("got more than 1 record...")
	}
	return nil
}

type Constructor func(columns []string, values []interface{}) map[string]interface{}

func DBSelectWithConstructor(sql string, constructor Constructor, args ...interface{}) []map[string]interface{} {
	rows, err := SqlDB.Query(sql, args...)
	if err != nil {
		log.Default().Println("failed: ", err)
	}

	defer rows.Close()

	cols, _ := rows.Columns()
	lenCol := len(cols)
	values := make([]interface{}, lenCol)
	for idx := range values {
		var val interface{}
		values[idx] = &val
	}

	var lst []map[string]interface{}
	for rows.Next() {
		err = rows.Scan(values...)
		if err != nil {
			log.Default().Println("failed: ", err)
		}

		item := constructor(cols, values)

		lst = append(lst, item)
	}
	return lst
}

func DBSelect(sql string, args ...interface{}) []map[string]interface{} {

	rows, err := SqlDB.Query(sql, args...)
	if err != nil {
		log.Default().Println("failed: ", err)
	}

	defer rows.Close()

	cols, _ := rows.Columns()
	lenCol := len(cols)
	values := make([]interface{}, lenCol)
	for idx := range values {
		var val interface{}
		values[idx] = &val
	}

	var lst []map[string]interface{}
	for rows.Next() {
		err = rows.Scan(values...)
		if err != nil {
			log.Default().Println("failed: ", err)
		}

		item := make(map[string]interface{})
		for idx, data := range values {
			item[cols[idx]] = *data.(*interface{})
		}

		lst = append(lst, item)
	}
	return lst
}

func DBInsert(sql string, args ...interface{}) int64 {
	ret, err := SqlDB.Exec(sql, args...)
	if err != nil {
		log.Default().Println("failed: ", err)
	}
	theID, err := ret.LastInsertId() // 新插入数据的id
	if err != nil {
		log.Default().Println("failed: ", err)
	}
	return theID
}

func DBUpdate(sql string, args ...interface{}) (int64, error) {
	ret, err := SqlDB.Exec(sql, args...)
	if err != nil {
		log.Default().Println("failed: ", err)
	}
	return ret.RowsAffected()
}

func DBDelete(sql string, args ...interface{}) {
	_, err := SqlDB.Exec(sql, args...)
	if err != nil {
		log.Default().Println("failed: ", err)
	}
}

func BulkInsert(sql string, values [][]interface{}) {
	stmt, err := SqlDB.Prepare(sql)
	if err != nil {
		log.Default().Println(err)
	}
	tx, err := SqlDB.Begin()
	if err != nil {
		log.Default().Println(err)
	}
	defer tx.Rollback()

	for _, value := range values {
		_, err := tx.Stmt(stmt).Exec(value...)
		if err != nil {
			log.Default().Println(err)
		}
	}

	err = tx.Commit()
	if err != nil {
		log.Default().Println(err)
	}

}
