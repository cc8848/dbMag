package main

import (
	"bytes"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"strings"
	"time"
	"strconv"
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
func SendData(tb *TableInfo, pip chan <-*ota_pre_record,maxid int ) error {

	db, err := tb.dbconn.GetConn()
	if err!=nil{
		return err
	}

	for i:=0;i<maxid;i+=50{
		j:=i+50
		sql:="select "+tb.colum+" from "+tb.name+" where id>"+strconv.Itoa(i) +" and id<="+strconv.Itoa(j)
		fmt.Println("sql string:",sql)
		rows,err:=db.Query(sql)
		defer rows.Close()
		if err!=nil{
			fmt.Println("querey error:",err)
		}

		for rows.Next()  {
			ota_pre := new(ota_pre_record)
			rows.Scan(&ota_pre.mid,
				&ota_pre.device_id,
				&ota_pre.product_id,
				&ota_pre.delta_id,
				&ota_pre.origin_version,
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
			//fmt.Println("send to buffer message:",ota_pre)
		}

	}


	return nil
}

/*

获取data缓冲区中的 TableInfo数据 发送到另外一个数据源
*/
func GetData(tb *TableInfo,pip <-chan *ota_pre_record) error {


	for true{
		ota_pre:=<-pip
		fmt.Println("revice from message :", ota_pre)
		db,err:=tb.dbconn.GetConn()

		if err!=nil{
			return nil
		}
		var sqlbuf bytes.Buffer
		sqlbuf.WriteString("INSERT INTO ")
		sqlbuf.WriteString(tb.name)
		sqlbuf.WriteString("(")
		sqlbuf.WriteString(tb.colum)
		sqlbuf.WriteString(") values(")
		colunArry:=strings.Split(tb.colum,",")
		size:=len(colunArry)
		for i:=0;i<size;i++{

			if i == size-1 {
				sqlbuf.WriteString("?")
			}else{
				sqlbuf.WriteString("?,")
			}
		}
		sqlbuf.WriteString(")")
		stmt, err :=db.Prepare(sqlbuf.String())
		stmt.Exec(ota_pre.mid,ota_pre.device_id,ota_pre.product_id,ota_pre.delta_id,ota_pre.origin_version,
			ota_pre.now_version,
			ota_pre.check_time,
			ota_pre.download_time,
			ota_pre.update_time,
			ota_pre.ip,
			ota_pre.province,
			ota_pre.city,
			ota_pre.networkType,
			ota_pre.status,
			ota_pre.origin_type,
			ota_pre.error_code,
			ota_pre.create_time,
			ota_pre.update_time)

		//fmt.Println(res.RowsAffected())
	}


	return nil
}


func main() {

	//缓冲区
	ch:=make( chan *ota_pre_record,50)

	//数据源1s
	connSrc:=DbConn{"180.97.81.42","root","123","dbconfig","33068"}
	tbinfoSrc:=TableInfo{dbconn:connSrc,name:"ota_pre_record",colum:"mid,device_id,product_id,delta_id,origin_version,now_version,check_time,download_time,upgrade_time,ip,province,city,networkType,status,origin_type,error_code,create_time,update_time"}
	db, err := tbinfoSrc.dbconn.GetConn()

	if err != nil {
		fmt.Println("connection source database error:",err)
	}
	rows,err:=db.Query("SELECT  max(id) from "+tbinfoSrc.name )
	var maxid int
	for rows.Next(){
		rows.Scan(&maxid)

	}

	go SendData(&tbinfoSrc,ch,maxid)

	//数据源2
	connDst:=DbConn{"180.97.81.42","root","123","dbconfig","33069"}

	tbinfoDst:=TableInfo{dbconn:connDst,name:"ota_pre_record",colum:"mid,device_id,product_id,delta_id,origin_version,now_version,check_time,download_time,upgrade_time,ip,province,city,networkType,status,origin_type,error_code,create_time,update_time"}


	go GetData(&tbinfoDst,ch)
	time.Sleep(10 * time.Second)

	fmt.Println("hello world")
}
