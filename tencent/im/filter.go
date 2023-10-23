package im

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
)

// A:B：A首次发送消息给B
// 触发发送单聊消息前回调
// 	- 检查redis中键（24小时自动过期）：key = A:B 的值 value = 秒时间戳
// 	- key不存在，表明24小时内首次发送消息，可以发送当前消息
// 	- key存在，表明24小时内发送过消息：
// 		- 检查redis中键 key = B:A 是否存在
// 		- 存在，表明24小时内B回复过A消息，可以继续发送当前消息
// 		- 不存在，表明24小时内B没有回复A消息，不能继续发送消息

const (
	// imFilterSingleChatKey 单聊消息限制
	imFilterSingleChatKey = "IMFilterSingleChat:%v:%v"
)

func FilterSingleSendKey(fromUid, toUid int32) string {
	return fmt.Sprintf(imFilterSingleChatKey, fromUid, toUid)
}

type IFilterCache interface {
	SetFlag(ctx context.Context, fromUid, toUid int32) error
	GetFlag(ctx context.Context, fromUid, toUid int32) (int64, error)
	ForbiddenSendMsg(ctx context.Context, fromUid, toUid int32) (bool, error)
}

type FilterCache struct {
	cache   *redis.Client
	options *FilterOptions
}

type FilterOptions struct {
	Expire int64 `json:"expire"` // 时间间隔限制
}

func NewFilterCache(cache *redis.Client, options *FilterOptions) *FilterCache {
	return &FilterCache{
		cache,
		options,
	}
}

func (f *FilterCache) SetFlag(ctx context.Context, fromUid, toUid int32) error {
	return f.cache.WithContext(ctx).SetEX(ctx, FilterSingleSendKey(fromUid, toUid), time.Now().Unix(), time.Duration(f.options.Expire)*time.Second).Err()
}

func (f *FilterCache) GetFlag(ctx context.Context, fromUid, toUid int32) (int64, error) {
	t, err := f.cache.WithContext(ctx).Get(ctx, FilterSingleSendKey(fromUid, toUid)).Int64()
	if err == redis.Nil {
		return 0, nil
	}
	return t, err
}

func (f *FilterCache) ForbiddenSendMsg(ctx context.Context, fromUid, toUid int32) (bool, error) {
	_, err := f.cache.WithContext(ctx).Get(ctx, FilterSingleSendKey(fromUid, toUid)).Int64()
	if err == redis.Nil {
		// A:B key不存在，可以发送消息
		return false, nil
	}
	if err != nil {
		return false, err
	}
	// key存在，表明时间限制内A给B发送过消息
	// 检查B是否24小时内给A发过消息：B:A
	_, err = f.cache.WithContext(ctx).Get(ctx, FilterSingleSendKey(toUid, fromUid)).Int64()
	if err == redis.Nil {
		// B:A key不存在，B时间限制内未给A发送过任何消息，禁止A发送消息
		return true, nil
	}
	if err != nil {
		return false, err
	}
	// 存在，可以发送消息
	return false, nil
}
