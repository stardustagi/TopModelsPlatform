package requests

type RegisterReq struct {
	Mail     string `json:"mail" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type LoginReq struct {
	Mail     string `json:"mail" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type UserIdReq struct {
	Id int64 `json:"id" query:"id" validate:"required"`
}

type CreateUserReq struct {
	UserName string `json:"user_name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Phone    string `json:"phone"`
	Password string `json:"password" validate:"required,min=6"`
}

type UpdateUserReq struct {
	Id       int64  `json:"id" validate:"required"`
	UserName string `json:"user_name"`
	Phone    string `json:"phone"`
	RealName string `json:"real_name"`
	Password string `json:"password" validate:"omitempty,min=6"`
}

type SetActiveReq struct {
	Id     int64 `json:"id" validate:"required"`
	Active int   `json:"active" validate:"oneof=0 1"`
}

// AdminAddUserReq 管理员添加用户请求
type AdminAddUserReq struct {
	UserName string `json:"user_name" validate:"required,min=3,max=20,alphanum"` // 用户名必填，3-20字符
	Password string `json:"password" validate:"required,min=6,max=20"`           // 密码必填，6-20字符
	Phone    string `json:"phone" validate:"required"`                           // 手机号必填
}

// UserAdminPaymentReq 用户充值请求
type UserAdminPaymentReq struct {
	UserName string `json:"user_name" validate:"required"`
	Amount   int64  `json:"amount" validate:"required,gt=0"`
	Reason   string `json:"reason" validate:"required"`
}

type GetModelsProviderInfoListReq struct {
	ProviderName string `json:"provider_name"`
}
