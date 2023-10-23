package im

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNewFilterCache(t *testing.T) {
	cli := NewFilterCache(redisCache, &FilterOptions{
		Expire: 24 * 3600,
	})
	assert.NotNil(t, cli, "NewFilterCache 失败")
	assert.Greater(t, cli.options.Expire, int64(0), "配置错误")
	assert.NotNil(t, cli.cache, "cache配置错误")
}

func TestImFilterSingleSendKey(t *testing.T) {
	key := FilterSingleSendKey(1, 2)
	assert.Equal(t, "IMFilterSingleChat:1:2", key, "生成key错误")
}

func TestFilterCache_SetFlag(t *testing.T) {
	cli := NewFilterCache(redisCache, &FilterOptions{
		Expire: 24 * 3600,
	})
	ctx := context.Background()
	key := FilterSingleSendKey(1, 2)
	err := redisCache.Del(ctx, key).Err()
	assert.Nil(t, err, err)

	ts, err := cli.GetFlag(ctx, 1, 2)
	assert.Nil(t, err, "获取不存在的键失败")
	assert.Equal(t, int64(0), ts, "获取值错误错误")

	err = cli.SetFlag(ctx, 1, 2)
	assert.Nil(t, err, err)

	ts, err = cli.GetFlag(ctx, 1, 2)
	assert.Nil(t, err, "获取存在的键失败")
	assert.Equal(t, time.Now().Unix(), ts, "获取值错误错误")

}
