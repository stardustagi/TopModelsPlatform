package responses

// GetGraphVerifyCodeResp 获取图形验证码响应
type GetGraphVerifyCodeResp struct {
	Code       string `json:"code"`       // 图形验证码
	ExpireTime string `json:"expireTime"` // 过期时间（秒）
}

// UserLoginAndRegisterRes 用户登录注册响应
type UserLoginAndRegisterRes struct {
	Success   bool        `json:"success"`
	Message   string      `json:"message"`
	UserID    int64       `json:"user_id"`
	Token     string      `json:"token"`
	IsNewUser bool        `json:"is_new_user"`
	UserInfo  interface{} `json:"user_info"`
}
