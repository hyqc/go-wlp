package im

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	"testing"
)

var redisCache *redis.Client

func init() {
	redisCache = redis.NewClient(&redis.Options{
		Addr: "192.168.9.100:6379",
		DB:   0,
	})
}

func TestNewIMOnlineStatusBitmap(t *testing.T) {

	ctx := context.Background()
	imKey := OnlineAllKey()
	ibm := NewIMOnlineStatusBitmap(redisCache)

	userIds := []int32{1, 3, 635, 636, 234, 12, 56, 12, 80, 11}
	for _, uid := range userIds {
		err := ibm.SetGlobalOnlineUserStatus(ctx, imKey, uid, true)
		assert.Nil(t, err, "设置全局在线状态失败", err)
	}

	for _, uid := range userIds {
		status, err := ibm.GetGlobalOnlineUserStatus(ctx, imKey, uid)
		assert.Nil(t, err, "获取全局在线状态失败", err)
		assert.Equal(t, true, status, "获取全局在线状态值错误")
	}

	tmpIds := []int32{1, 636, 80, 11, 999, 1201}
	expected := []int64{1, 1, 1, 1, 0, 0}
	status, err := ibm.BatchGetGlobalOnlineUserStatus(ctx, imKey, OnlineMutualUserIdKey(635), tmpIds, BitfieldSetValueTypeU1)
	assert.Nil(t, err, "批量获取某个玩家的关系用户状态失败", err)
	assert.Equal(t, expected, status, "预期值不对", status)

	tmpIds = []int32{1}
	expected = []int64{1}
	status, err = ibm.BatchGetGlobalOnlineUserStatus(ctx, imKey, OnlineMutualUserIdKey(635), tmpIds, BitfieldSetValueTypeU1)
	assert.Nil(t, err, "批量获取某个玩家的关系用户状态失败", err)
	assert.Equal(t, expected, status, "预期值不对", status)
}

func TestIMOnlineStatusBitmap_BatchGetGlobalOnlineUserStatus(t *testing.T) {
	ctx := context.Background()
	imKey := OnlineAllKey()
	ibm := NewIMOnlineStatusBitmap(redisCache)

	userIds := []int32{1}
	for _, uid := range userIds {
		err := ibm.SetGlobalOnlineUserStatus(ctx, imKey, uid, true)
		assert.Nil(t, err, "设置全局在线状态失败", err)
	}

	for _, uid := range userIds {
		status, err := ibm.GetGlobalOnlineUserStatus(ctx, imKey, uid)
		assert.Nil(t, err, "获取全局在线状态失败", err)
		assert.Equal(t, true, status, "获取全局在线状态值错误")
	}

	tmpIds := []int32{1}
	expected := []int64{1}
	status, err := ibm.BatchGetGlobalOnlineUserStatus(ctx, imKey, OnlineMutualUserIdKey(635), tmpIds, BitfieldSetValueTypeU1)
	assert.Nil(t, err, "批量获取某个玩家的关系用户状态失败", err)
	assert.Equal(t, expected, status, "预期值不对", status)
	tmpIds = []int32{6}
	expected = []int64{0}
	status, err = ibm.BatchGetGlobalOnlineUserStatus(ctx, imKey, OnlineMutualUserIdKey(635), tmpIds, BitfieldSetValueTypeU1)
	assert.Nil(t, err, "批量获取某个玩家的关系用户状态失败", err)
	assert.Equal(t, expected, status, "预期值不对", status)

}

func TestIMOnlineMutualUserIdKey(t *testing.T) {
	ctx := context.Background()
	imKey := OnlineAllKey()
	ibm := NewIMOnlineStatusBitmap(redisCache)
	tmpKey := OnlineMutualUserIdKey(636)
	u1, err := ibm.BatchGetGlobalOnlineUserStatusU1(ctx, imKey, tmpKey, []int32{635})
	assert.Nil(t, err, err)
	fmt.Println(u1)
}
