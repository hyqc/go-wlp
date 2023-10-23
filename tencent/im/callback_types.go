package im

/*****************************/
// IM在线状态回调

type CallbackOnlineStatusReqInfo struct {
	Action    string `json:"Action"`     // Login,Logout,Disconnect,CustomStatusChange
	ToAccount string `json:"To_Account"` // 用户 UserID
	Reason    string `json:"Reason"`     //
}
type KickedDeviceItem struct {
	Platform string `json:"Platform"` // 平台
}

type CallbackOnlineStatusReq struct {
	CallbackCommand string                      `json:"CallbackCommand"` // 回调命令
	EventTime       int64                       `json:"EventTime"`       // 触发本次回调的时间戳，单位为毫秒。
	Info            CallbackOnlineStatusReqInfo `json:"Info"`
	KickedDevice    []KickedDeviceItem          `json:"KickedDevice"`
}

const (
	OK       = "OK"   // 成功
	OkCode   = 0      // 成功
	FAIL     = "FAIL" // 失败
	FailCode = 1      // 失败
)

const (
	// CustomCodeErrNoReply IM回调自定义错误码120001-130000之间
	CustomCodeErrNoReply = 120002 // 24时间限制内未得到对方回应
)

// CodeMsg 错误消息
var CodeMsg = map[int]string{
	CustomCodeErrNoReply: "24小时只能发送1条消息",
}

type CallbackOnlineStatusResp struct {
	ActionStatus string `json:"ActionStatus"`
	ErrorCode    int    `json:"ErrorCode"` // 0成功，1失败
	ErrorInfo    string `json:"ErrorInfo"`
}

/*****************************/
// 发送单聊消息之前回调

type SingleMsgBodyItem struct {
	MsgType    string         `json:"MsgType"` // 消息类型
	MsgContent map[string]any `json:"MsgContent"`
}

type CallbackSingleSendBeforeReq struct {
	CallbackCommand string              `json:"CallbackCommand"` // 回调命令
	FromAccount     string              `json:"From_Account"`    // 发送者
	ToAccount       string              `json:"To_Account"`      // 接收者
	MsgSeq          int64               `json:"MsgSeq"`          // 消息序列号
	MsgRandom       int64               `json:"MsgRandom"`       // 消息随机数
	MsgTime         int64               `json:"MsgTime"`         // 消息的发送时间戳，单位为秒
	MsgKey          string              `json:"MsgKey"`          //消息的唯一标识，可用于 REST API 撤回单聊消息
	OnlineOnlyFlag  int32               `json:"OnlineOnlyFlag"`  //在线消息，为1，否则为0；
	MsgBody         []SingleMsgBodyItem `json:"MsgBody"`         // 消息体
	CloudCustomData string              `json:"CloudCustomData"`
	EventTime       int64               `json:"EventTime"` // 触发本次回调的时间戳，单位为毫秒。
}

type CallbackSingleSendBeforeResp struct {
	ActionStatus    string              `json:"ActionStatus"` // OK, FAIL
	ErrorCode       int32               `json:"ErrorCode"`    // 0为允许发言  1为禁止发言  2为静默丢弃
	ErrorInfo       string              `json:"ErrorInfo"`    // 错误信息
	MsgBody         []SingleMsgBodyItem `json:"MsgBody"`
	CloudCustomData string              `json:"CloudCustomData"` // 消息自定义数据
}
