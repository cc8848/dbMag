package client

import (
	"bytes"
	"crypto/tls"
	"encoding/binary"
	"errors"
	"mysql"
	"net"
	"packet"
	"strings"
	"time"
	"utils"
)

/*
客户端链接属性
*/
type ClientConn struct {
	*packet.Packet
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

	return nil
}

func (c *ClientConn) writeAuthHandShake() error {
	// Adjust client capability flags based on server support
	capability := mysql.CLIENT_PROTOCOL_41 | mysql.CLIENT_SECURE_CONNECTION |
		mysql.CLIENT_LONG_PASSWORD | mysql.CLIENT_TRANSACTIONS | mysql.CLIENT_LONG_FLAG

	// To enable TLS / SSL
	if c.TLSConfig != nil {
		capability |= mysql.CLIENT_PLUGIN_AUTH
		capability |= mysql.CLIENT_SSL
	}

	capability &= c.capability

	//packet length
	//capbility 4
	//max-packet size 4
	//charset 1
	//reserved all[0] 23
	length := 4 + 4 + 1 + 23

	//username
	length += len(c.user) + 1

	//we only support secure connection
	auth := mysql.CalcPassword(c.salt, []byte(c.password))

	length += 1 + len(auth)

	if len(c.db) > 0 {
		capability |= mysql.CLIENT_CONNECT_WITH_DB

		length += len(c.db) + 1
	}

	// mysql_native_password + null-terminated
	length += 21 + 1

	c.capability = capability

	data := make([]byte, length+4)

	//capability [32 bit]
	data[4] = byte(capability)
	data[5] = byte(capability >> 8)
	data[6] = byte(capability >> 16)
	data[7] = byte(capability >> 24)

	//MaxPacketSize [32 bit] (none)
	//data[8] = 0x00
	//data[9] = 0x00
	//data[10] = 0x00
	//data[11] = 0x00

	//Charset [1 byte]
	//use default collation id 33 here, is utf-8
	data[12] = byte(mysql.DEFAULT_COLLATION_ID)

	// SSL Connection Request Packet
	// http://dev.mysql.com/doc/internals/en/connection-phase-packets.html#packet-Protocol::SSLRequest
	if c.TLSConfig != nil {
		// Send TLS / SSL request packet
		if err := c.WritePacket(data[:(4+4+1+23)+4]); err != nil {
			return err
		}

		// Switch to TLS
		tlsConn := tls.Client(c.Packet.Conn, c.TLSConfig)
		if err := tlsConn.Handshake(); err != nil {
			return err
		}

		currentSequence := c.Sequence
		c.Packet = packet.NewConn(tlsConn)
		c.Sequence = currentSequence
	}

	//Filler [23 bytes] (all 0x00)
	pos := 13 + 23

	//User [null terminated string]
	if len(c.user) > 0 {
		pos += copy(data[pos:], c.user)
	}
	data[pos] = 0x00
	pos++

	// auth [length encoded integer]
	data[pos] = byte(len(auth))
	pos += 1 + copy(data[pos+1:], auth)

	// db [null terminated string]
	if len(c.db) > 0 {
		pos += copy(data[pos:], c.db)
		data[pos] = 0x00
		pos++
	}

	// Assume native client during response
	pos += copy(data[pos:], "mysql_native_password")
	data[pos] = 0x00

	return c.WritePacket(data)
	return nil
}

/*
解析server下发的握手信息包

*/
func (conn *ClientConn) readInitialHandShake() error {
	data, err := conn.ReadPacket()
	if err != nil {
		return err
	}

	if data[0] == mysql.ERR_HEADER {
		return errors.New("read initial handshake error")
	}

	if data[0] < mysql.MinProtocolVersion {
		return errors.New("invalid protocol version")
	}

	//skip mysql version
	//mysql version end with 0x00
	pos := 1 + bytes.IndexByte(data[1:], 0x00) + 1

	//connection id length is 4
	conn.connectionID = uint32(binary.LittleEndian.Uint32(data[pos : pos+4]))

	pos += 4
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
func (c *ClientConn) Close() error {
	return c.Conn.Close()
}

func (c *ClientConn) ReadOKPacket() (*mysql.Result, error) {
	return c.readOK()
}

func (c *ClientConn) readOK() (*mysql.Result, error) {
	data, err := c.ReadPacket()
	if err != nil {
		return nil, err
	}

	if data[0] == mysql.OK_HEADER {
		return c.handleOKPacket(data)
	} else if data[0] == mysql.ERR_HEADER {
		return nil, c.handleErrorPacket(data)
	} else {
		return nil, errors.New("invalid ok packet")
	}
}
func (c *ClientConn) handleErrorPacket(data []byte) error {
	e := new(mysql.MyError)

	var pos int = 1

	e.Code = binary.LittleEndian.Uint16(data[pos:])
	pos += 2

	if c.capability&mysql.CLIENT_PROTOCOL_41 > 0 {
		//skip '#'
		pos++
		e.State = utils.String(data[pos : pos+5])
		pos += 5
	}

	e.Message = utils.String(data[pos:])

	return e
}
func (c *ClientConn) readUntilEOF() (err error) {
	var data []byte

	for {
		data, err = c.ReadPacket()

		if err != nil {
			return
		}

		// EOF Packet
		if c.isEOFPacket(data) {
			return
		}
	}
	return
}

func (c *ClientConn) isEOFPacket(data []byte) bool {
	return data[0] == mysql.EOF_HEADER && len(data) <= 5
}

func (c *ClientConn) handleOKPacket(data []byte) (*mysql.Result, error) {
	var n int
	var pos int = 1

	r := new(mysql.Result)

	r.AffectedRows, _, n = mysql.LengthEncodedInt(data[pos:])
	pos += n
	r.InsertId, _, n = mysql.LengthEncodedInt(data[pos:])
	pos += n

	if c.capability&mysql.CLIENT_PROTOCOL_41 > 0 {
		r.Status = binary.LittleEndian.Uint16(data[pos:])
		c.status = r.Status
		pos += 2

		//todo:strict_mode, check warnings as error
		//Warnings := binary.LittleEndian.Uint16(data[pos:])
		//pos += 2
	} else if c.capability&mysql.CLIENT_TRANSACTIONS > 0 {
		r.Status = binary.LittleEndian.Uint16(data[pos:])
		c.status = r.Status
		pos += 2
	}

	//new ok package will check CLIENT_SESSION_TRACK too, but I don't support it now.

	//skip info
	return r, nil
}

func (c *ClientConn)writeCommand(command byte) error  {

	c.Packet.ResetSequence()
	return c.WritePacket([]byte{0x01,0x00,0x00,0x00,command})
}

/*
发生字符串命令
*/
func (c *ClientConn)writeCommandStr(command byte,args string) error  {

	c.Packet.ResetSequence()
	length:=len(args)+1

	data:=make([]byte,length+4)
	data[4]=command

	copy(data[5:],args)
	return c.WritePacket(data)
}



func (c *ClientConn)readResultColumns(result *mysql.Result) (err error){

	var i int =0
	var data []byte

	for{
		data,err=c.ReadPacket()
		if err!=nil{
			return
		}

		//EOF packet
		if c.isEOFPacket(data){
			if c.capability&mysql.CLIENT_PROTOCOL_41 >0{
				result.Status=binary.LittleEndian.Uint16(data[3:])
				c.status=result.Status
			}

			if i!=len(result.Fields){
				err=mysql.ErrMalformPacket
			}

			return
		}

		result.Fields[i],err=mysql.FieldData(data).Parse()
		if err!=nil{
			return
		}

		result.FieldNames[utils.String(result.Fields[i].Name)]=i
		i++
	}
}


func (c *ClientConn)readResultset(data []byte,bl bool)(*mysql.Result,error){

	result:=&mysql.Result{

		Status:0,
		InsertId:0,
		AffectedRows:0,
		Resultset:&mysql.Resultset{},
	}

	count,_,n:=mysql.LengthEncodedInt(data)
	if n-len(data) !=0{
		return nil,mysql.ErrMalformPacket
	}

	result.Fields=make([]*mysql.Field,count)
	result.FieldNames=make(map[string]int,count)

	if err:=c.readResultColumns(result);err!=nil{
		return nil,err
	}

	if err:=c.readResultRows(result,bl);err!=nil{
		return nil,err
	}

	return result,nil
}

func (c *ClientConn) readResultRows(result *mysql.Result, isBinary bool) (err error) {
	var data []byte

	for {
		data, err = c.ReadPacket()

		if err != nil {
			return
		}

		// EOF Packet
		if c.isEOFPacket(data) {
			if c.capability&mysql.CLIENT_PROTOCOL_41 > 0 {
				//result.Warnings = binary.LittleEndian.Uint16(data[1:])
				//todo add strict_mode, warning will be treat as error
				result.Status = binary.LittleEndian.Uint16(data[3:])
				c.status = result.Status
			}

			break
		}

		result.RowDatas = append(result.RowDatas, data)
	}

	result.Values = make([][]interface{}, len(result.RowDatas))

	for i := range result.Values {
		result.Values[i], err = result.RowDatas[i].Parse(result.Fields, isBinary)

		if err != nil {
			return  err
		}
	}

	return nil
}


func (c *ClientConn)readResult(bl bool)(*mysql.Result,error)  {

	data,err:=c.ReadPacket()
	if err!=nil{
		return nil,err
	}

	if data[0]==mysql.OK_HEADER{
		return c.handleOKPacket(data)
	}else if data[0]==mysql.ERR_HEADER{
		return nil,c.handleErrorPacket(data)
	}else if data[0]==mysql.LocalInFile_HEADER{
		return nil,mysql.ErrMalformPacket
	}

	return c.readResultset(data,bl)
}


func(c *ClientConn)Ping()error{

	if err:=c.writeCommand(mysql.COM_PING);err!=nil{
		return err
	}
	if _,err:=c.readOK();err!=nil{
		return err
	}

	return nil
}

func (c *ClientConn)exec(query string)(*mysql.Result,error)  {
	if err:=c.writeCommandStr(mysql.COM_QUERY,query);err!=nil{
		return nil,err
	}

	return c.readResult(false)
}


func (c *ClientConn) writeCommandUint32(command byte, arg uint32) error {
	c.Packet.ResetSequence()

	return c.WritePacket([]byte{
		0x05, //5 bytes long
		0x00,
		0x00,
		0x00, //sequence

		command,

		byte(arg),
		byte(arg >> 8),
		byte(arg >> 16),
		byte(arg >> 24),
	})
}

//执行命令入口
func (c *ClientConn)Execute(command string,args...interface{}) (*mysql.Result,error)  {

	if len(args)==0{
		return c.exec(command)
	}else {

		if s,err:=c.Prepare(command);err!=nil{
			return nil,err
		}else{
			var r *mysql.Result
			r,err=s.Execute(args...)
			s.Close()
			return r,err
		}
	}
}