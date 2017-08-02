package main

import (
	"bytes"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

type TableInfo struct {
	dbconn DbConn //为所在的数据库 connection
	name   string //表的名字
	idval  uint64 //表中最大的id值 表必须有一个id属性而且为主键
}

type DbConn struct {
	host   string
	user   string
	passwd string
	db     string
	port   uint32
}

type TableInterFace interface {
	GetData(tb TableInfo, data interface{})
	SendData(tb TableInfo)
}

func (conn *DbConn) GetConn() (*sql.DB, error) {
	var dataSourceName bytes.Buffer
	dataSourceName.WriteString(conn.user)
	dataSourceName.WriteString(":")
	dataSourceName.WriteString(conn.passwd)
	dataSourceName.WriteString("@tcp(")
	dataSourceName.WriteString(conn.host)
	dataSourceName.WriteString(":")
	dataSourceName.WriteString(")/")
	dataSourceName.WriteString(conn.db)
	dataSourceName.WriteString("/charset=utf8")
	db, err := sql.Open("tcp", dataSourceName.String())
	if err != nil {
		return nil, err
	}

	return db, nil
}

func GetData(tb TableInfo, data interface{}) error {

	db, err := tb.dbconn.GetConn()
	if err != nil {
		return err
	}
	rows, err := db.Query("SELECT  max(id) from " + tb.name)

	rows.Scan(&tb.idval)

	return nil
}

func main() {

	fmt.Println("hello world")
}
