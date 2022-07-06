package pool

import (
	"github.com/golang/protobuf/proto"
	"math"
	"sync"
)

type CacheBuffer struct {
	proto.Buffer
	lastMarshaledSize uint32
}

var BufferPool = &sync.Pool{
	New: func() interface{} {
		return &CacheBuffer{
			Buffer:            proto.Buffer{},
			lastMarshaledSize: 16,
		}
	},
}

func (b *CacheBuffer) GetLastMarshaledSize() uint32 {
	return b.lastMarshaledSize
}
func (b *CacheBuffer) SetLastMarshaledSize(newSize int) {
	if newSize > math.MaxUint32 {
		newSize = math.MaxUint32
	}

	b.lastMarshaledSize = uint32(newSize)
}