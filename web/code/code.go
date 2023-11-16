package code

import (
	"web/pkg/response"
)

const (
	Success    = 0
	Failed     = 1
	AuthFailed = 100001 // 鉴权
)

const (
	SuccessMsg    = "成功"
	FailedMsg     = "失败"
	AuthFailedMsg = "未登录或登录已过期"
)

var Msg = map[int]string{
	Success:    SuccessMsg,
	Failed:     FailedMsg,
	AuthFailed: AuthFailedMsg,
}

func NewCode(code int) response.Message {
	return response.Message{
		MessageBase: response.MessageBase{
			Code:    code,
			Message: Msg[code],
		},
	}
}

func NewCodeMsg(code int, msg string) response.Message {
	return response.Message{
		MessageBase: response.MessageBase{
			Code:    code,
			Message: msg,
		},
	}
}

func NewCodeMsgData(code int, msg string, data interface{}) response.Message {
	return response.Message{
		MessageBase: response.MessageBase{
			Code:    code,
			Message: msg,
		},
		Data: data,
	}
}
