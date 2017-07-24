package main

import (
	"crypto/tls"
	"errors"
	"mysql"
	"net"
	"packet"
	"strings"
	"time"
	"bytes"
	"encoding/binary"
)

/*
客户端链接属性
*/
type ClientConn struct {
	*packet.Conn
	user      string
	password  string
	db        string
	TLSConfig *tls.Config

	capability uint32
	status     uint16
	charset    string

	connectionID uint32
	salt         []byte
}

func Connect(addr string, user string, password string, dbName string, options ...func(*ClientConn)) (*ClientConn, error) {
	var proto string
	if strings.Contains(addr, "/") {
		proto = "unix"
	} else {
		proto = "tcp"
	}

	var err error
	conn, err := net.DialTimeout(proto, addr, 10*time.Second)
	if err != nil {
		return nil, err
	}
	c := new(ClientConn)
	c.Conn = packet.NewConn(conn)
	c.user = user
	c.password = password
	c.db = dbName

	c.charset = mysql.DEFAULT_CHARSET
	// Apply configuration functions.
	for i := range options {
		options[i](c)
	}

	if err = c.handshake(); err != nil {
		return nil, err
	}

	return c, nil
}
func (c *ClientConn) handshake() error {

	var err error
	//读取server下发的握手信息包
	if err = c.readInitialHandShake(); err != nil {
		c.Close()
		return err
	}

	//写认证消息包
	if err = c.writeAuthHandShake(); err != nil {
		c.Close()
		return err
	}

	//读取server下发是否ok的消息
	if err := c.readOk(); err != nil {
		c.Close()
		return err
	}

	return nil
}

func (conn *ClientConn) readOk() error {

	return  nil
}

func (conn *ClientConn) writeAuthHandShake() error {

	return nil
}
/*
解析server下发的握手信息包

*/
func (conn *ClientConn) readInitialHandShake() error {
	data,err:=conn.ReadPacket()
	if err!=nil{
		return err
	}

	if data[0]==mysql.ERR_HEADER{
		return errors.New("read initial handshake error")
	}

	if data[0] <mysql.MinProtocolVersion{
		return errors.New("invalid protocol version")
	}

	//skip mysql version
	//mysql version end with 0x00
	pos := 1 + bytes.IndexByte(data[1:], 0x00) + 1


	//connection id length is 4
	conn.connectionID=uint32(binary.LittleEndian.Uint32(data[pos : pos+4]))

	pos+=4
	conn.salt = []byte{}
	conn.salt = append(conn.salt, data[pos:pos+8]...)
	//skip filter
	pos += 8 + 1

	//capability lower 2 bytes
	conn.capability = uint32(binary.LittleEndian.Uint16(data[pos : pos+2]))

	pos += 2

	if len(data) > pos {
		//skip server charset
		//c.charset = data[pos]
		pos += 1

		conn.status = binary.LittleEndian.Uint16(data[pos : pos+2])
		pos += 2

		conn.capability = uint32(binary.LittleEndian.Uint16(data[pos:pos+2]))<<16 | conn.capability

		pos += 2

		//skip auth data len or [00]
		//skip reserved (all [00])
		pos += 10 + 1

		// The documentation is ambiguous about the length.
		// The official Python library uses the fixed length 12
		// mysql-proxy also use 12
		// which is not documented but seems to work.
		conn.salt = append(conn.salt, data[pos:pos+12]...)
	}


	return nil
}
