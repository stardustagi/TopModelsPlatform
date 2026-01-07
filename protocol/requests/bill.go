package requests

// CreateModelsDiscountReq 创建模型折扣请求
type CreateModelsDiscountReq struct {
	ModelId         int64  `json:"model_id" validate:"required"`
	ModelProviderId int64  `json:"model_provider_id"`
	Type            string `json:"type" validate:"required"` // percent, money, time
	Value           int64  `json:"value" validate:"required"`
}

// GetModelsDiscountReq 获取模型折扣请求
type GetModelsDiscountReq struct {
	Id int64 `json:"id" validate:"required"`
}

// ListModelsDiscountReq 获取模型折扣列表请求
type ListModelsDiscountReq struct {
	ModelId         int64   `json:"model_id"`
	ModelProviderId int64   `json:"model_provider_id"`
	Type            string  `json:"type"`
	PageInfo        PageReq `json:"page_info"`
}

// UpdateModelsDiscountReq 更新模型折扣请求
type UpdateModelsDiscountReq struct {
	Id              int64  `json:"id" validate:"required"`
	ModelId         int64  `json:"model_id"`
	ModelProviderId int64  `json:"model_provider_id"`
	Type            string `json:"type"`
	Value           int64  `json:"value"`
}

// DeleteModelsDiscountReq 删除模型折扣请求
type DeleteModelsDiscountReq struct {
	Id int64 `json:"id" validate:"required"`
}

// GetModelDiscountByModelIdReq 根据模型ID获取折扣请求
type GetModelDiscountByModelIdReq struct {
	ModelId         int64 `json:"model_id" validate:"required"`
	ModelProviderId int64 `json:"model_provider_id"`
}

type CreateUserRebateConfigReq struct {
	UserId     int64 `json:"user_id" validate:"required,gt=0"`
	TierStart  int64 `json:"tier_start" validate:"gte=0"`
	TierEnd    int64 `json:"tier_end" validate:"required,gtefield=TierStart|eq=-1"`
	RebateRate int   `json:"rebate_rate" validate:"required,gte=0,lte=100"`
	Status     *int  `json:"status"` // 使用指针，nil时默认启用，0表示禁用，1表示启用
}

type UpdateUserRebateConfigReq struct {
	Id         int64 `json:"id" validate:"required,gt=0"`
	UserId     int64 `json:"user_id" validate:"required,gt=0"`
	TierStart  int64 `json:"tier_start" validate:"gte=0"`
	TierEnd    int64 `json:"tier_end" validate:"required,gtefield=TierStart|eq=-1"`
	RebateRate int   `json:"rebate_rate" validate:"required,gte=0,lte=100"`
	Status     int   `json:"status" validate:"oneof=0 1"`
}

type GetUserRebateConfigReq struct {
	PageInfo   PageReq `json:"page_info"`
	UserId     int64   `json:"user_id"`
	TierStart  int64   `json:"tier_start"`
	TierEnd    int64   `json:"tier_end"`
	RebateRate int     `json:"rebate_rate"`
	Status     *int    `json:"status"`
	CreatedAt  int64   `json:"created_at"`
}

type GetUserRebateInfoReq struct {
	PageInfo      PageReq `json:"page_info"`
	UserId        int64   `json:"user_id"`
	Month         string  `json:"month"`
	TotalConsumed int64   `json:"total_consumed"`
	RebateRate    int     `json:"rebate_rate"`
}

// GetUserConsumeRecordReq 获取用户消费记录请求
type GetUserConsumeRecordReq struct {
	PageInfo       PageReq `json:"page_info"`
	UserId         int64   `json:"user_id" validate:"gte=0"`
	StartTime      int64   `json:"start_time"`
	EndTime        int64   `json:"end_time"`
	Model          string  `json:"model"`
	ProviderId     string  `json:"provider_id"`
	ActualProvider string  `json:"actual_provider"`
	CreatedAt      int64   `json:"created_at"`
	TotalConsumed  int64   `json:"total_consumed"`
}

// UserConsumeDetailReq 获取用户消费详情请求
type UserConsumeDetailReq struct {
	ConsumeId   int64   `json:"consume_id" validate:"required"`
	ConsumeType string  `json:"consume_type"`
	Page        PageReq `json:"page"`
}

// GetDayReportListReq 获取日报表列表请求
type GetDayReportListReq struct {
	UserId       int64   `json:"user_id"`
	UserName     string  `json:"user_name"`
	ProviderName string  `json:"provider_name"`
	ModelName    string  `json:"model_name"`
	ConsumeType  string  `json:"consume_type"`
	UpdatedAt    int64   `json:"updated_at"`
	TotalCost    int64   `json:"total_cost"`
	Page         PageReq `json:"page"`
}
