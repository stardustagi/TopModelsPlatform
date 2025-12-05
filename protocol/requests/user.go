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
}

type SetActiveReq struct {
	Id     int64 `json:"id" validate:"required"`
	Active int   `json:"active" validate:"oneof=0 1"`
}
