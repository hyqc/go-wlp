package redisutil

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewIMOnlineStatusBitmap(t *testing.T) {
	client := redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
	})
	defer func(client *redis.Client) {
		err := client.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(client)

	ctx := context.Background()
	imKey := IMOnlineAllKey()
	ibm := NewIMOnlineStatusBitmap(client)

	userIds := []int64{1, 3, 635, 636, 234, 12, 56, 12, 80, 11}
	for _, uid := range userIds {
		err := ibm.SetGlobalOnlineUserStatus(ctx, imKey, uid, true)
		assert.Nil(t, err, "设置全局在线状态失败", err)
	}

	for _, uid := range userIds {
		status, err := ibm.GetGlobalOnlineUserStatus(ctx, imKey, uid)
		assert.Nil(t, err, "获取全局在线状态失败", err)
		assert.Equal(t, true, status, "获取全局在线状态值错误")
	}

	tmpIds := []int64{1, 636, 80, 11, 999, 1201}
	expected := []int64{1, 1, 1, 1, 0, 0}
	status, err := ibm.BatchGetGlobalOnlineUserStatus(ctx, imKey, IMOnlineMutualUserIdKey(635), tmpIds, BitfieldSetValueTypeU1)
	assert.Nil(t, err, "批量获取某个玩家的关系用户状态失败", err)
	assert.Equal(t, expected, status, "预期值不对", status)
}
