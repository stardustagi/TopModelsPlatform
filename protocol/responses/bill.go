package responses

import "github.com/stardustagi/TopModelsPlatform/models"

// CreateModelsDiscountResp 创建模型折扣响应
type CreateModelsDiscountResp struct {
	Id              int64  `json:"id"`
	ModelId         int64  `json:"model_id"`
	ModelProviderId int64  `json:"model_provider_id"`
	Type            string `json:"type"`
	Value           int64  `json:"value"`
	CreatedAt       int64  `json:"created_at"`
}

// GetModelsDiscountResp 获取模型折扣响应
type GetModelsDiscountResp struct {
	Id              int64  `json:"id"`
	ModelId         int64  `json:"model_id"`
	ModelProviderId int64  `json:"model_provider_id"`
	Type            string `json:"type"`
	Value           int64  `json:"value"`
	CreatedAt       int64  `json:"created_at"`
	LastUpdate      int64  `json:"last_update"`
}

// ListModelsDiscountResp 获取模型折扣列表响应
type ListModelsDiscountResp struct {
	Discounts []models.ModelsDiscount `json:"discounts"`
	Total     int                     `json:"total"`
}

// UpdateModelsDiscountResp 更新模型折扣响应
type UpdateModelsDiscountResp struct {
	Id              int64  `json:"id"`
	ModelId         int64  `json:"model_id"`
	ModelProviderId int64  `json:"model_provider_id"`
	Type            string `json:"type"`
	Value           int64  `json:"value"`
	LastUpdate      int64  `json:"last_update"`
}

type UserRebateData struct {
	UserId        int64  `json:"user_id" xorm:"'user_id' BIGINT(20)"`
	Month         string `json:"month" xorm:"'month' VARCHAR(7)"`                             // 格式：2024-01
	TotalConsumed int64  `json:"total_consumed" xorm:"'total_consumed' BIGINT(20) default 0"` // 当月消费总额
	RebateAmount  int64  `json:"rebate_amount" xorm:"'rebate_amount' BIGINT(20) default 0"`   // 已返点金额
	RebateUsed    int64  `json:"rebate_used" xorm:"'rebate_used' BIGINT(20) default 0"`       // 已消费返点
	RebateRate    int    `json:"rebate_rate" xorm:"'rebate_rate' INT(10)"`                    // 返点比例快照
	Status        int    `json:"status" xorm:"'status' INT(10) default 0"`                    // 0未返点 1已返点

}

type GetUserRebateConfigResp struct {
	RebateInfoList []models.UserRebateConfig `json:"rebate_info_list"`
	Total          int                       `json:"total"`
}

type GetUserRebateInfoResp struct {
	Data  []UserRebateData `json:"data"`
	Total int              `json:"total"`
}

type DayReportItem struct {
	Id               int64  `json:"id"`
	UserId           int64  `json:"user_id"`
	UserName         string `json:"user_name"`
	ActualProviderId int    `json:"actual_provider_id"`
	ProviderName     string `json:"provider_name"`
	ModelId          int    `json:"model_id"`
	ModelName        string `json:"model_name"`
	ConsumeType      string `json:"consume_type"`
	Date             string `json:"date"`
	TotalConsumed    int64  `json:"total_consumed"`
	TotalCost        int64  `json:"total_cost"`
	UpdatedAt        int64  `json:"updated_at"`
}

type GetDayReportListResp struct {
	Items []DayReportItem `json:"items"`
	Total int             `json:"total"`
	Skip  int             `json:"skip"`
	Limit int             `json:"limit"`
}
