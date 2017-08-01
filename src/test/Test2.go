package main

import (
	"strings"
	"bytes"
	"fmt"
)

func main() {
	var sqlbuf bytes.Buffer
	colum:="mid,device_id,product_id,delta_id"

	sqlbuf.WriteString("INSERT INTO ")
	sqlbuf.WriteString("ota_pre_record")
	sqlbuf.WriteString("(")
	sqlbuf.WriteString(colum)
	sqlbuf.WriteString(") values(")

	colunArry:=strings.Split(colum,",")
	size:=len(colunArry)
	for i:=0;i<size;i++{

		if i == size-1 {
			sqlbuf.WriteString("?")
		}else{
			sqlbuf.WriteString("?,")
		}
	}
	sqlbuf.WriteString(")")

	fmt.Println(sqlbuf.String())
}