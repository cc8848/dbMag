package main

import "packet"

type ClientConn struct {
	*packet.Conn
	user string
	password string
	db string

	connectionID uit32
	salt []byte

	status uint8
}

