package requests

// ===================== DiscountRule 请求 =====================

// CreateDiscountRuleReq 创建折扣规则请求
type CreateDiscountRuleReq struct {
	Name         string `json:"name" validate:"required"`
	Description  string `json:"description"`
	DiscountRate int    `json:"discount_rate" validate:"required,gte=0,lte=100"`
	Status       *int   `json:"status"`
}

// UpdateDiscountRuleReq 更新折扣规则请求
type UpdateDiscountRuleReq struct {
	Id           int64  `json:"id" validate:"required,gt=0"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	DiscountRate *int   `json:"discount_rate"`
	Status       *int   `json:"status"`
}

// GetDiscountRuleListReq 获取折扣规则列表请求
type GetDiscountRuleListReq struct {
	PageInfo     PageReq `json:"page_info"`
	Name         string  `json:"name"`
	DiscountRate *int    `json:"discount_rate"`
	Status       *int    `json:"status"`
	CreatedAt    int64   `json:"created_at"`
}

// ===================== UserDiscount 请求 =====================

// CreateUserDiscountReq 创建用户折扣请求
type CreateUserDiscountReq struct {
	UserId  int64 `json:"user_id" validate:"required,gt=0"`
	ModelId int   `json:"model_id" validate:"required,gt=0"`
	RuleId  int64 `json:"rule_id"`
}

type DiscountContent struct {
	Id     int64 `json:"id" validate:"required,gt=0"`
	RuleId int64 `json:"rule_id"`
}

// UpdateUserDiscountReq 更新用户折扣请求
type UpdateUserDiscountReq struct {
	Data []DiscountContent `json:"data" validate:"required,dive"`
}

// GetUserDiscountListReq 获取用户折扣列表请求
type GetUserDiscountListReq struct {
	PageInfo     PageReq `json:"page_info"`
	UserId       int64   `json:"user_id"`
	ModelId      int     `json:"model_id"`
	RuleId       int64   `json:"rule_id"`
	DiscountRate *int    `json:"discount_rate"`
	CreatedAt    int64   `json:"created_at"`
}
