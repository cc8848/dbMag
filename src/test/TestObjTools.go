package main

import (
	"fmt"

	"bytes"
	"encoding/gob"
	"log"
)

type Stu struct {
	Name string
	Age  uint32
}

func main() {

	var stu2 Stu
	stu1 := Stu{Name: "gao", Age: 30}

	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)

	var str string = "hello gaohaicang!"
	e, err := Encode(str)
	if err != nil {
		log.Print(err.Error())
	}

	log.Println(e)

	//fmt.Println("hello world!")
}

func Encode(data interface{}) ([]byte, error) {

	buf := bytes.NewBuffer(nil)

	enc := gob.NewEncoder(buf)
	err := enc.Encode(data)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func Decode(data []byte, to interface{}) error {

	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)

	return dec.Decode(to)
}
