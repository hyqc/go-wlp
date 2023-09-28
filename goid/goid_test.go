package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetGid(t *testing.T) {
	gid, err := GetGid()
	fmt.Println(gid)
	assert.Nil(t, err, "获取GID失败")
	assert.NotZero(t, gid, "获取GID失败")
}

func BenchmarkGetGid(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = GetGid()
	}
}

func BenchmarkGetGidN(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = GetGidN()
	}
}
