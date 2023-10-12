package redisutil

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"strconv"
	"strings"
)

// bitfield key2 set u1 userid1 1 set u1 userid2 1 set u1 userid3 1
// bitop and key2 key1 key1
// bitfield_ro key2 get u1 userid1 get u1 userid2 get u1 userid3
// del key2

const (
	BitfieldSetValueTypeU1 = "u1" // 无符号1位
)

const (
	// IM 全部玩家的在线状态 bitmap key
	imOnlineAllKey = "IMOnline:All"
	// IM 我的相互关注列表的玩家的在线状态临时 bitmap key，用于计算列表的玩家状态
	imOnlineMutualUserIdKey = "IMOnline:Mutual:%v"
)

func IMOnlineAllKey() string {
	return imOnlineAllKey
}

func IMOnlineMutualUserIdKey(userId int64) string {
	return fmt.Sprintf(imOnlineMutualUserIdKey, userId)
}

type IIMOnlineStatusBitmap interface {
	SetGlobalOnlineUserStatus(ctx context.Context, key string, userId int64, isOnline bool) error // 设置全局在线状态值
	GetGlobalOnlineUserStatus(ctx context.Context, key string, userId int64) (bool, error)        // 获取全局在线状态
	BatchGetGlobalOnlineUserStatus(ctx context.Context, key, tmpKey string, userIds []int64, vt string) ([]int64, error)
	BatchGetGlobalOnlineUserStatusU1(ctx context.Context, key, tmpKey string, userIds []int64) ([]int64, error)
	BatchGetUsersOnlineStatus(userIds []int64, status []int64) map[int64]bool
}

type IMOnlineStatusBitmap struct {
	client *redis.Client
}

func NewIMOnlineStatusBitmap(client *redis.Client) *IMOnlineStatusBitmap {
	return &IMOnlineStatusBitmap{
		client,
	}
}

func (i *IMOnlineStatusBitmap) SetGlobalOnlineUserStatus(ctx context.Context, key string, userId int64, isOnline bool) error {
	status := 0
	if isOnline {
		status = 1
	}
	return i.client.SetBit(ctx, key, userId, status).Err()
}

func (i *IMOnlineStatusBitmap) GetGlobalOnlineUserStatus(ctx context.Context, key string, userId int64) (bool, error) {
	status, err := i.client.GetBit(ctx, key, userId).Result()
	if err != nil {
		return false, err
	}
	return status == 1, nil
}

func (i *IMOnlineStatusBitmap) BatchGetGlobalOnlineUserStatus(ctx context.Context, key, tmpKey string, userIds []int64, vt string) ([]int64, error) {
	scp := i.batchMakeScript()
	keys, args := i.batchHandleParams(key, tmpKey, userIds, vt)
	return redis.NewScript(scp).Eval(ctx, i.client, keys, args).Int64Slice()
}

func (i *IMOnlineStatusBitmap) BatchGetGlobalOnlineUserStatusU1(ctx context.Context, key, tmpKey string, userIds []int64) ([]int64, error) {
	return i.BatchGetGlobalOnlineUserStatus(ctx, key, tmpKey, userIds, BitfieldSetValueTypeU1)
}

func (i *IMOnlineStatusBitmap) BatchGetUsersOnlineStatus(userIds []int64, result []int64) map[int64]bool {
	m := make(map[int64]bool)
	for i, uid := range userIds {
		m[uid] = result[i] == 1
	}
	return m
}

// batchMakeScript 生成批量获取bitmap中多个用户在线状态的lua 脚本字符串
func (i *IMOnlineStatusBitmap) batchMakeScript() string {
	script := `
local key = KEYS[1];
local tmpKey = KEYS[2];
local u1 = ARGV[1];
local str = ARGV[2];

local ids = {}
local start = 1
for i = 1, #str do
	local first, last = str:find("(,)", start)
	local id = 0
	if first then
		id = tonumber(str:sub(start, first - 1))
		table.insert(ids, id)
		start = last + 1
	else
		id = tonumber(str:sub(start))
		table.insert(ids, id)
		break
	end
end

-- 设置临时bitmap
for i,v in ipairs(ids) do
	redis.call('bitfield', tmpKey, 'set', u1, v, 1)
end

-- bitop 计算
redis.call('bitop','and', tmpKey, key, tmpKey);

-- 获取临时bitmap的值
local status = {}
for i, v in ipairs(ids) do
    status[i] = redis.call('bitfield_ro', tmpKey, 'get', u1, v)[1];
end

-- 删除临时键
redis.call('del', tmpKey)

return status;
`
	return script
}

// batchHandleParams batchMakeScript() 需要的参数处理
func (i *IMOnlineStatusBitmap) batchHandleParams(key, tmpKey string, userIds []int64, vt string) (keys []string, args []interface{}) {
	keys = []string{key, tmpKey}
	args = []interface{}{vt}
	ids := make([]string, 0, len(userIds))
	for _, id := range userIds {
		ids = append(ids, strconv.FormatInt(id, 10))
	}
	args = append(args, strings.Join(ids, ","))
	return
}
