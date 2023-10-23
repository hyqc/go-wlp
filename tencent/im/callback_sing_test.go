package im

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestSign_Check(t *testing.T) {
	sn := NewSign(&SignOptions{
		Token:        "123456A",
		Expire:       60,
		IgnoreExpire: true,
	})
	var rt int64 = 1697075895
	rs := "a990654049bdb6caa85f3763098570d25f02e49dfc9e7f04960367bca54f7615"
	err := sn.Check(rs, rt)
	assert.Nil(t, err, " sign check failed")

	rs = "a990654049bdb6caa85f3763098570d25f02e49dfc9e7f04960367bca54f761511111"
	err = sn.Check(rs, rt)
	assert.NotNil(t, err, " sign check failed")

	sn = NewSign(&SignOptions{
		Token:        "123456A",
		Expire:       60,
		IgnoreExpire: false,
	})
	rt = 1697075895
	rs = "a990654049bdb6caa85f3763098570d25f02e49dfc9e7f04960367bca54f7615"
	err = sn.Check(rs, rt)
	assert.NotNil(t, err, " sign expired check failed")
	rt = time.Now().Unix()
	err = sn.Check(rs, rt)
	assert.NotNil(t, err, "request sign error")
	ns := sn.Make(rt)
	err = sn.Check(ns, rt)
	assert.Nil(t, err, "pass")
}

func TestSign_Make(t *testing.T) {
	sn := NewSign(&SignOptions{
		Token: "123456A",
	})
	var rt int64 = 1669872112
	ns := sn.Make(rt)
	rs := "fb1561b2e2daf6e98d4fcaf7a26d1ed02608f9f9669721f7ead1ff0f7b4f32ca"
	assert.Equal(t, rs, ns, "make sign error")
}

func TestCompareTime(t *testing.T) {
	var a int64 = 1669872112
	var b int64 = 1669872114
	res := CompareTime(a, b, 3*time.Second)
	assert.Truef(t, res, "A+3s在B后")
	res = CompareTime(a, b, 1*time.Second)
	assert.Falsef(t, res, "A+1s在B前")
	res = CompareTime(a, b, 0)
	assert.Falsef(t, res, "A+0s在B前")
	res = CompareTime(a, b, 2)
	assert.Falsef(t, res, "A+2s=B")
}
