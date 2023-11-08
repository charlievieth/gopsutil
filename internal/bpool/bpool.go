package bpool

import (
	"bufio"
	"io"
	"sync"
)

var pool = sync.Pool{
	New: func() interface{} {
		return bufio.NewReader(nil)
	},
}

func GetReader(rd io.Reader) *bufio.Reader {
	r := pool.Get().(*bufio.Reader)
	r.Reset(rd)
	return r
}

func PutReader(r *bufio.Reader) {
	if r != nil {
		r.Reset(nil)
		pool.Put(r)
	}
}
