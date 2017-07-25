package packet

import (
	"bufio"
	"bytes"
	"io"
	"net"
	"mysql"
	"errors"
)

/*
	Conn is the base class to handle MySQL protocol.
*/
type Packet struct {
	net.Conn
	br *bufio.Reader

	Sequence uint8
}

func NewConn(conn net.Conn) *Packet {
	p := new(Packet)

	p.br = bufio.NewReaderSize(conn, 4096)
	p.Conn = conn

	return p
}


func (p *Packet) ReadPacket() ([]byte, error) {
	var buf bytes.Buffer

	if err := p.ReadPacketTo(&buf); err != nil {
		return nil,err
	} else {
		return buf.Bytes(), nil
	}

	// header := []byte{0, 0, 0, 0}

	// if _, err := io.ReadFull(c.br, header); err != nil {
	// 	return nil, ErrBadConn
	// }

	// length := int(uint32(header[0]) | uint32(header[1])<<8 | uint32(header[2])<<16)
	// if length < 1 {
	// 	return nil, fmt.Errorf("invalid payload length %d", length)
	// }

	// sequence := uint8(header[3])

	// if sequence != c.Sequence {
	// 	return nil, fmt.Errorf("invalid sequence %d != %d", sequence, c.Sequence)
	// }

	// c.Sequence++

	// data := make([]byte, length)
	// if _, err := io.ReadFull(c.br, data); err != nil {
	// 	return nil, ErrBadConn
	// } else {
	// 	if length < MaxPayloadLen {
	// 		return data, nil
	// 	}

	// 	var buf []byte
	// 	buf, err = c.ReadPacket()
	// 	if err != nil {
	// 		return nil, ErrBadConn
	// 	} else {
	// 		return append(data, buf...), nil
	// 	}
	// }
}

func (p *Packet) ReadPacketTo(w io.Writer) error {
	header := []byte{0, 0, 0, 0}

	if _, err := io.ReadFull(p.br, header); err != nil {
		return mysql.ErrBadConn
	}

	length := int(uint32(header[0]) | uint32(header[1])<<8 | uint32(header[2])<<16)
	if length < 1 {
		return errors.New("invalid payload length")
	}

	sequence := uint8(header[3])

	if sequence != p.Sequence {
		return errors.New("invalid sequence")
	}

	p.Sequence++

	if n, err := io.CopyN(w, p.br, int64(length)); err != nil {
		return mysql.ErrBadConn
	} else if n != int64(length) {
		return mysql.ErrBadConn
	} else {
		if length < mysql.MaxPayloadLen {
			return nil
		}

		if err := p.ReadPacketTo(w); err != nil {
			return err
		}
	}

	return nil
}

// data already has 4 bytes header
// will modify data inplace
func (p *Packet) WritePacket(data []byte) error {
	length := len(data) - 4

	for length >= mysql.MaxPayloadLen {
		data[0] = 0xff
		data[1] = 0xff
		data[2] = 0xff

		data[3] = p.Sequence

		if n, err := p.Write(data[:4+mysql.MaxPayloadLen]); err != nil {
			return mysql.ErrBadConn
		} else if n != (4 + mysql.MaxPayloadLen) {
			return mysql.ErrBadConn
		} else {
			p.Sequence++
			length -= mysql.MaxPayloadLen
			data = data[mysql.MaxPayloadLen:]
		}
	}

	data[0] = byte(length)
	data[1] = byte(length >> 8)
	data[2] = byte(length >> 16)
	data[3] = p.Sequence

	if n, err := p.Write(data); err != nil {
		return mysql.ErrBadConn
	} else if n != len(data) {
		return mysql.ErrBadConn
	} else {
		p.Sequence++
		return nil
	}
}

func (p *Packet) ResetSequence() {
	p.Sequence = 0
}

func (p *Packet) Close() error {
	p.Sequence = 0
	if p.Conn != nil {
		return p.Conn.Close()
	}
	return nil
}
