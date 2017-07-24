package main

import (
	"packet"
	"crypto/tls"
	"strings"
	"net"
	"time"
	"mysql"
	"errors"
)
/*
客户端链接属性
*/
type ClientConn struct {
	*packet.Conn
	user     string
	password string
	db       string
	TLSConfig *tls.Config

	capability uint32
	status uint16
	charset string

	connectionID uint32
	salt         []byte

}

func Connect(addr string, user string, password string, dbName string, options ...func(*ClientConn)) (*ClientConn, error)  {
	var proto string
	if strings.Contains(addr,"/"){
		proto="unix"
	}else {
		proto="tcp"
	}

	var err error
	conn,err:=net.DialTimeout(proto,addr,10*time.Second)
	if err!=nil{
		return nil,err
	}
	c:=new(ClientConn)
	c.Conn=packet.NewConn(conn)
	c.user=user
	c.password=password
	c.db=dbName

	c.charset=mysql.DEFAULT_CHARSET
	// Apply configuration functions.
	for i := range options {
		options[i](c)
	}

	if err = c.handshake(); err != nil {
		return nil, err
	}

	return c,nil
}
func (c *ClientConn) handshake() error {

	var err error
	//读取server下发的握手信息包
	if err=c.readInitialHandShake();err!=nil{
		c.Close()
		return err
	}

	//写认证消息包
	if err=c.writeAuthHandShake();err!=nil{
		c.Close()
		return err
	}

	//读取server下发是否ok的消息
	if err:=c.readOk();err!=nil{
		c.Close()
		return err
	}

	return nil
}




func (conn *ClientConn) readOk() error {

}


func (conn *ClientConn) writeAuthHandShake() error {

}


func (conn *ClientConn) readInitialHandShake() error {

}
