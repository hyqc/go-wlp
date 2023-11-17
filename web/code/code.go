package code

import (
	"web/pkg/response"
)

const (
	Success = 0
	Failed  = 1

	AuthTokenFailed         = 200001 // 鉴权
	AuthTokenInvalid        = 200002 // BearerToken无效
	AuthTokenInspectInvalid = 200003 // BearerToken无效
	AuthTokenInfoInvalid    = 200004 // Token信息无效

	RequestBodyInvalid   = 300001 // 请求参数无效
	RequestQueryInvalid  = 300002 // 请求参数无效
	RequestParamsInvalid = 300003 // 请求参数无效
)

const (
	SuccessMsg = "成功"
	FailedMsg  = "失败"

	AuthFailedMsg              = "未登录或登录令牌已过期"
	AuthTokenInvalidMsg        = "登录令牌无效"
	AuthTokenInspectInvalidMsg = "登录令牌检查失败"
	AuthTokenInfoInvalidMsg    = "登录令牌信息无效"

	RequestBodyInvalidMsg   = "请求体参数无效"
	RequestQueryInvalidMsg  = "查询参数无效"
	RequestParamsInvalidMsg = "请求参数无效"
)

var Msg = map[int]string{
	Success: SuccessMsg,
	Failed:  FailedMsg,

	AuthTokenFailed:         AuthFailedMsg,
	AuthTokenInvalid:        AuthTokenInvalidMsg,
	AuthTokenInspectInvalid: AuthTokenInspectInvalidMsg,
	AuthTokenInfoInvalid:    AuthTokenInfoInvalidMsg,

	RequestBodyInvalid:   RequestBodyInvalidMsg,
	RequestQueryInvalid:  RequestQueryInvalidMsg,
	RequestParamsInvalid: RequestParamsInvalidMsg,
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
