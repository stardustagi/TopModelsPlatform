package responses

import "github.com/stardustagi/TopModelsPlatform/models"

// ===================== DiscountRule 响应 =====================

// GetDiscountRuleListResp 获取折扣规则列表响应
type GetDiscountRuleListResp struct {
	List  []models.DiscountRule `json:"list"`
	Total int                   `json:"total"`
}

// ===================== UserDiscount 响应 =====================

// GetUserDiscountListResp 获取用户折扣列表响应
type GetUserDiscountListResp struct {
	List  []models.UserDiscount `json:"list"`
	Total int                   `json:"total"`
}

// 折扣日志
type GetDiscountLogListResp struct {
	List  []models.UserConsumeRecord `json:"list"`
	Total int64                      `json:"total"`
}
