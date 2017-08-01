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
	colum  string //列名 格式 a,b,c
	idval  uint32 //表中最大的id值 表必须有一个id属性而且为主键
	tabledata *ota_pre_record
}

type DbConn struct {
	host   string
	user   string
	passwd string
	db     string
	port   string
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
func SendData(tb *TableInfo, pip chan *ota_pre_record ) error {

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

	rows2,err:=db.Query("select "+tb.colum+" from "+tb.name," where id>0 and id<20")
	for rows2.Next()  {
		ota_pre := new(ota_pre_record)
		rows.Scan(&ota_pre.mid,
			&ota_pre.device_id,
			&ota_pre.product_id,
			&ota_pre.delta_id,
			&ota_pre.origin_type,
			&ota_pre.now_version,
			&ota_pre.check_time,
			&ota_pre.download_time,
			&ota_pre.update_time,
			&ota_pre.ip,
			&ota_pre.province,
			&ota_pre.city,
			&ota_pre.networkType,
			&ota_pre.status,
			&ota_pre.origin_type,
			&ota_pre.error_code,
			&ota_pre.create_time,
			&ota_pre.update_time,
		)
		pip<-ota_pre
	}
	return nil
}

/*

获取data缓冲区中的 TableInfo数据
*/
func GetData(tb *TableInfo,pip chan *ota_pre_record) error {



	return nil
}


func main() {

	ch:=make( chan ota_pre_record,10)
	conn:=DbConn{"180.97.81.42","root","123","dbconfig","33068"}


	tbinfo:=TableInfo{dbconn:conn,name:"ota_pre_record",colum:"mid,device_id,product_id,delta_id,origin_version,now_version,check_time,download_time,upgrade_time,ip,province,city,networkType,status,origin_type,error_code,create_time,update_time"}

	GetData(&tbinfo,ch)

	fmt.Println("max id values is:",tbinfo.idval)
	fmt.Println("hello world")
}
