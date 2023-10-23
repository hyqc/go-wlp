package im

type IM struct {
	SdkAppId      int            `json:"sdkappid"`        // IM的应用SdkAppId
	Key           string         `json:"key"`             // IM的密钥
	UserSigExpire int            `json:"user_sig_expire"` // 过期时间，单位秒
	Callback      *SignOptions   `json:"callback"`        // 回调鉴权相关
	Filter        *FilterOptions `json:"filter"`          // 聊天限制
}
