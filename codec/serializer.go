package codec

import (
	"RPC_go/pool"
	"errors"
	"github.com/golang/protobuf/proto"
)

var serializerMap = map[string]Serializer{
	"proto": &ProtoSerializer{},
}

func GetSerialization(sType string) Serializer {
	if serializer, ok := serializerMap[sType]; ok && serializer != nil {
		return serializer
	}
	return &ProtoSerializer{}
}

type Serializer interface {
	Serialize(data interface{}) ([]byte, error)
	Deserialize(data []byte, resp interface{}) error
}

type ProtoSerializer struct{}

// Serialize convert proto to []bytes
func (s *ProtoSerializer) Serialize(data interface{}) ([]byte, error) {
	if data == nil {
		return nil, nil
	}

	var msgData proto.Message
	if marshal, ok := data.(proto.Marshaler); ok {
		return marshal.Marshal()

	} else {
		msgData = data.(proto.Message)
	}

	cache := pool.BufferPool.Get().(*pool.CacheBuffer)
	defer pool.BufferPool.Put(cache)

	buf := make([]byte, 0, cache.GetLastMarshaledSize())
	cache.SetBuf(buf)
	cache.Reset()
	if err := cache.Marshal(msgData); err != nil {
		return nil, err
	}

	bData := cache.Bytes()
	cache.SetLastMarshaledSize(len(bData))
	cache.SetBuf(nil)

	return bData, nil
}

// Deserialize convert []bytes to proto struct
func (s *ProtoSerializer) Deserialize(bData []byte, resp interface{}) error {
	if bData == nil || len(bData) == 0 {
		return errors.New("[codec]deserialize failed, err=nil data")
	}

	msgData := resp.(proto.Message)
	msgData.Reset()

	cache := pool.BufferPool.Get().(*pool.CacheBuffer)
	defer pool.BufferPool.Put(cache)

	cache.SetBuf(bData)
	err := cache.Unmarshal(msgData)
	cache.SetBuf(nil)
	return err
}
