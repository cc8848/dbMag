package main

import (
	"fmt"
	"client"
)
func main() {
	addr:="180.97.81.42:33068"
	user:="root"
	password:="123"
	db:="dbconfig"
	c,err:=client.Connect(addr,user,password,db)
	if err!=nil{
		fmt.Println("connection database fail")
	}
	sql:="create table tst(id int,imei varchar(10))"
	c.Execute(sql)
	fmt.Println("hello world!")
}