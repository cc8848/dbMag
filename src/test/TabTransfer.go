package main

import "fmt"

type TableInfo struct {
	dbconn DbConn //为所在的数据库 connection
	name   string //表的名字
	maxid  uint64 //表中最大的id值 表必须有一个id属性而且为主键
}

type DbConn struct {
	host   string
	user   string
	passwd string
	port   uint32
}

type TableInterFace interface {
	GetData()
	SendData()
}




func main() {
	fmt.Println("hello world")
}
