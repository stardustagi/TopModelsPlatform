package requests

// CreateAlarmReq 创建告警配置请求
type CreateAlarmReq struct {
	Type   string `json:"type" validate:"required,oneof=token billing"`
	Min    int    `json:"min" validate:"gte=0"`
	Max    int    `json:"max" validate:"gtefield=Min"`
	Status *int   `json:"status"` // 使用指针区分未传值和0
	UserId *int64 `json:"user_id"`
}

// UpdateAlarmReq 更新告警配置请求
type UpdateAlarmReq struct {
	Id     int64  `json:"id" validate:"required,gt=0"`
	Type   string `json:"type" validate:"omitempty,oneof=token billing"`
	Min    *int   `json:"min"`
	Max    *int   `json:"max"`
	Status *int   `json:"status"`
	UserId *int64 `json:"user_id"`
}

// GetAlarmReq 获取告警配置请求
type GetAlarmReq struct {
	Id     int64   `json:"id"`
	UserId int64   `json:"user_id"`
	Type   string  `json:"type"`
	Page   PageReq `json:"page"`
}
