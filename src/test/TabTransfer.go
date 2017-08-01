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
	colum  []string
	idval  uint32 //表中最大的id值 表必须有一个id属性而且为主键
}

type DbConn struct {
	host   string
	user   string
	passwd string
	db     string
	port   string
}

type TableInterFace interface {
	GetData(tb *TableInfo, data interface{})
	SendData(tb *TableInfo,data interface{})
}

func (conn *DbConn) GetConn() (*sql.DB, error) {
	var dataSourceName bytes.Buffer
	dataSourceName.WriteString(conn.user)
	dataSourceName.WriteString(":")
	dataSourceName.WriteString(conn.passwd)
	dataSourceName.WriteString("@tcp(")
	dataSourceName.WriteString(conn.host)
	dataSourceName.WriteString(":")
	dataSourceName.WriteString(conn.port)
	dataSourceName.WriteString(")/")
	dataSourceName.WriteString(conn.db)
	dataSourceName.WriteString("?charset=utf8")
	db, err := sql.Open("mysql", dataSourceName.String())
	if err != nil {
		return nil, err
	}

	return db, nil
}


/*
将TableInfo 中的数据发送到缓冲区 data中
*/
func SendData(tb *TableInfo, pip chan ota_pre_record ) error {

	db, err := tb.dbconn.GetConn()

	if err != nil {
		return err
	}
	rows,err:=db.Query("SELECT  max(id) from ota_pre_record" )

	for rows.Next(){
		var maxid uint32
		rows.Scan(&maxid)
		tb.idval=maxid
	}





	return nil
}

/*

获取data缓冲区中的 TableInfo数据
*/
func GetData(tb *TableInfo,pip chan ota_pre_record) error {



	return nil
}

type ota_pre_record  struct {
	mid string
	device_id string
	product_id string
	delta_id string
	origin_version string
	now_version string
	check_time string
	download_time string
	upgrade_time string
	ip string
	province string
	city string
	networkType string
	status string
	origin_type string
	error_code string
	create_time string
	update_time string
}

func main() {

	ch:=make( chan ota_pre_record,10)
	conn:=DbConn{"180.97.81.42","root","123","dbconfig","33068"}
	tbinfo:=TableInfo{dbconn:conn,name:"ota_pre_record"}

	GetData(&tbinfo,ch)

	fmt.Println("max id values is:",tbinfo.idval)
	fmt.Println("hello world")
}
