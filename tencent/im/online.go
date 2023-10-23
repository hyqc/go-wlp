package im

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
	// IM 我的关注列表玩家在线状态临时 bitmap key，用于计算列表的玩家状态
	imOnlineTempFollowUserIdKey = "IMOnlineTemp:Follow:%v:"
	// IM 我的粉丝列表的玩家的在线状态临时 bitmap key，用于计算列表的玩家状态
	imOnlineTempFanUserIdKey = "IMOnlineTemp:Fan:%v"
	// IM 我的相互关注列表的玩家的在线状态临时 bitmap key，用于计算列表的玩家状态
	imOnlineTempMutualUserIdKey = "IMOnlineTemp:Mutual:%v"
)

func OnlineAllKey() string {
	return imOnlineAllKey
}

func OnlineFollowUserIdKey(userId int32) string {
	return fmt.Sprintf(imOnlineTempFollowUserIdKey, userId)
}

func OnlineFanUserIdKey(userId int32) string {
	return fmt.Sprintf(imOnlineTempFanUserIdKey, userId)
}

func OnlineMutualUserIdKey(userId int32) string {
	return fmt.Sprintf(imOnlineTempMutualUserIdKey, userId)
}

type IIMOnlineStatusBitmap interface {
	SetGlobalOnlineUserStatus(ctx context.Context, imKey string, userId int32, isOnline bool) error // 设置全局在线状态值
	GetGlobalOnlineUserStatus(ctx context.Context, imKey string, userId int32) (bool, error)        // 获取全局在线状态
	BatchGetGlobalOnlineUserStatus(ctx context.Context, imKey, tmpKey string, userIds []int32, vt string) ([]int64, error)
	BatchGetGlobalOnlineUserStatusU1(ctx context.Context, imKey, tmpKey string, userIds []int32) ([]int32, error)
	BatchHandleUsersOnlineStatus(userIds []int32, status []int32) (map[int32]bool, int)
}

type OnlineStatusBitmap struct {
	client *redis.Client
}

// NewIMOnlineStatusBitmap new client
func NewIMOnlineStatusBitmap(client *redis.Client) *OnlineStatusBitmap {
	return &OnlineStatusBitmap{
		client,
	}
}

// SetGlobalOnlineUserStatus 设置指定玩家的在线状态
func (i *OnlineStatusBitmap) SetGlobalOnlineUserStatus(ctx context.Context, imKey string, userId int32, isOnline bool) error {
	status := 0
	if isOnline {
		status = 1
	}
	return i.client.SetBit(ctx, imKey, int64(userId), status).Err()
}

// GetGlobalOnlineUserStatus 获取指定玩家的在线状态
func (i *OnlineStatusBitmap) GetGlobalOnlineUserStatus(ctx context.Context, imKey string, userId int32) (bool, error) {
	status, err := i.client.GetBit(ctx, imKey, int64(userId)).Result()
	if err != nil {
		return false, err
	}
	return status == 1, nil
}

// BatchGetGlobalOnlineUserStatus 批量获取在线状态
func (i *OnlineStatusBitmap) BatchGetGlobalOnlineUserStatus(ctx context.Context, imKey, tmpKey string, userIds []int32, vt string) ([]int64, error) {
	scp := i.batchMakeScript()
	keys, args := i.batchHandleParams(imKey, tmpKey, userIds, vt)
	return redis.NewScript(scp).Eval(ctx, i.client, keys, args).Int64Slice()
}

// BatchGetGlobalOnlineUserStatusU1 批量获取在线状态
func (i *OnlineStatusBitmap) BatchGetGlobalOnlineUserStatusU1(ctx context.Context, imKey, tmpKey string, userIds []int32) ([]int32, error) {
	data, err := i.BatchGetGlobalOnlineUserStatus(ctx, imKey, tmpKey, userIds, BitfieldSetValueTypeU1)
	if err != nil {
		return nil, err
	}
	res := make([]int32, 0, len(data))
	for _, v := range data {
		res = append(res, int32(v))
	}
	return res, nil
}

// BatchHandleUsersOnlineStatus 在线状态数据处理
func (i *OnlineStatusBitmap) BatchHandleUsersOnlineStatus(userIds []int32, states []int32) (map[int32]bool, int) {
	m := make(map[int32]bool)
	total := 0
	for i, uid := range userIds {
		state := states[i] == 1
		m[uid] = state
		if state {
			total++
		}
	}
	return m, total
}

func (i *OnlineStatusBitmap) BatchGetGlobalOnlineUserStatusU12Map(ctx context.Context, imKey, tmpKey string, userIds []int32) (states map[int32]bool, onlineNum int, err error) {
	tmpStates, err := i.BatchGetGlobalOnlineUserStatusU1(ctx, imKey, tmpKey, userIds)
	if err != nil {
		return nil, 0, err
	}
	states, onlineNum = i.BatchHandleUsersOnlineStatus(userIds, tmpStates)
	return states, onlineNum, nil
}

// redis-7.x版本支持，5.x版本不支持
// batchMakeScript 生成批量获取bitmap中多个用户在线状态的lua 脚本字符串
func (i *OnlineStatusBitmap) batchMakeScript() string {
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
    status[i] = redis.call('bitfield', tmpKey, 'get', u1, v)[1];
end

-- 删除临时键
redis.call('del', tmpKey)

return status;
`
	return script
}

// batchHandleParams batchMakeScript() 需要的参数处理
func (i *OnlineStatusBitmap) batchHandleParams(imKey, tmpKey string, userIds []int32, vt string) (keys []string, args []interface{}) {
	keys = []string{imKey, tmpKey}
	args = []interface{}{vt}
	ids := make([]string, 0, len(userIds))
	for _, id := range userIds {
		ids = append(ids, strconv.FormatInt(int64(id), 10))
	}
	args = append(args, strings.Join(ids, ","))
	return
}
