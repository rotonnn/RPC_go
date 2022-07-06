package codec

import (
	"RPC_go/constant"
	"bytes"
	"encoding/binary"
)

var codecMap = map[string]Codec{}

func GetCodec(protocolType string) Codec {
	if c, ok := codecMap[protocolType]; ok && c != nil {
		return c
	}
	return &defaultCodec{}
}

type FrameHeader struct {
	MagicNum     uint8
	Version      uint8
	MsgType      uint8
	ReqType      uint8
	CompressType uint8
	StreamID     uint16
	Length       uint32
	Reserved     uint32
}

type Codec interface {
	Encode([]byte) ([]byte, error)
	Decode([]byte) ([]byte, error)
}

type defaultCodec struct{}

func (c *defaultCodec) Encode(data []byte) ([]byte, error) {
	length := constant.FrameHeaderLength + len(data)
	buffer := bytes.NewBuffer(make([]byte, 0, length))

	frame := &FrameHeader{
		MagicNum: uint8(constant.MagicNumber),
		Version:  uint8(constant.Version),
		Length:   uint32(len(data)),
	}

	if err := binary.Write(buffer, binary.BigEndian, frame.MagicNum); err != nil {
		return nil, err
	}
	if err := binary.Write(buffer, binary.BigEndian, frame.Version); err != nil {
		return nil, err
	}
	if err := binary.Write(buffer, binary.BigEndian, frame.MsgType); err != nil {
		return nil, err
	}
	if err := binary.Write(buffer, binary.BigEndian, frame.ReqType); err != nil {
		return nil, err
	}
	if err := binary.Write(buffer, binary.BigEndian, frame.CompressType); err != nil {
		return nil, err
	}
	if err := binary.Write(buffer, binary.BigEndian, frame.StreamID); err != nil {
		return nil, err
	}
	if err := binary.Write(buffer, binary.BigEndian, frame.Length); err != nil {
		return nil, err
	}
	if err := binary.Write(buffer, binary.BigEndian, frame.Reserved); err != nil {
		return nil, err
	}
	if err := binary.Write(buffer, binary.BigEndian, data); err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

func (c *defaultCodec) Decode(data []byte) ([]byte, error) {
	return data[constant.FrameHeaderLength:], nil
}
