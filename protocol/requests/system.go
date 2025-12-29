package requests

// GetGraphVerifyCodeReq 图形验证码验证请求
type GetGraphVerifyCodeReq struct {
	T string `json:"t" validate:"required"` // 时间戳
}

// GetPhoneVerifyCodeReq 获取电话验证码请求
type GetPhoneVerifyCodeReq struct {
	GraphCode string `json:"graph_code" validate:"required"` // 图形验证码
	Phone     string `json:"phone" validate:"required"`      // 手机号
	T         string `json:"t" validate:"required"`          // 时间戳
	RegType   string `json:"regType" validate:"required"`    // 注册类型
}

// PhoneLoginAndRegisterReq 手机登录注册请求
type PhoneLoginAndRegisterReq struct {
	Phone      string `json:"phone" validate:"required"`      // 手机号必填
	PhoneCode  string `json:"phone_code" validate:"required"` // 手机验证码必填
	InviteCode string `json:"invite_code"`                    // 邀请码
	T          string `json:"t" validate:"required"`          // 时间戳必填
}
