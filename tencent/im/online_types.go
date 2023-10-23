package im

const (
	OnlineStatusOff uint8 = 0 // 下线
	OnlineStatusOn  uint8 = 1 // 上线
)

type OnlineStatusChangeBody struct {
	Uid    int32 `json:"uid"`
	Status uint8 `json:"status"` // 1上线，0下线
}

const (
	OnlineStatusChangeBrokerTopic = "IMOnlineStatusChange"
)
