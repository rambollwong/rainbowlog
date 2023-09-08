package rainbowlog

import (
	"github.com/rambollwong/rainbowcat/pool"
	"sync"
)

var bytesPool = pool.NewBytesPool(128, pool.DefaultMaxBytesCap)

type NewRecordFunc func() *Record

type recordPool struct {
	pool *sync.Pool
}

func newRecordPool(newFunc NewRecordFunc) *recordPool {
	return &recordPool{
		pool: &sync.Pool{
			New: func() interface{} {
				return newFunc()
			},
		},
	}
}

func (p *recordPool) Get() *Record {
	return p.pool.Get().(*Record)
}

func (p *recordPool) Put(r *Record) {
	r.Reset()
	p.pool.Put(r)
}
