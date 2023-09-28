package main

import (
	"bytes"
	"runtime"
	"strconv"
	"sync"
)

var (
	BufPoll64 = sync.Pool{
		New: func() any {
			bf := make([]byte, 64)
			return &bf
		},
	}
	prefixBytes = []byte("goroutine ")
)

func GetGid() (int64, error) {
	bp := BufPoll64.Get().(*[]byte)
	defer BufPoll64.Put(bp)
	b := *bp
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, prefixBytes)
	i := bytes.IndexByte(b, ' ')
	gid, err := strconv.ParseInt(string(b[:i]), 10, 64)
	if err != nil {
		return 0, err
	}
	return gid, nil
}

func GetGidN() (int64, error) {
	by := [64]byte{}
	b := by[:runtime.Stack(by[:], false)]
	b = bytes.TrimPrefix(b, prefixBytes)
	i := bytes.IndexByte(b, ' ')
	gid, err := strconv.ParseInt(string(b[:i]), 10, 64)
	if err != nil {
		return 0, err
	}
	return gid, nil
}
