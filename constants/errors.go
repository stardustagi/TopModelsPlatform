package constants

import (
	topError "github.com/stardustagi/TopLib/libs/errors"
)

var (
	ErrInternalServer    = topError.New("Internal server error", 500)
	ErrInvalidParams     = topError.New("无效的请求参数", 501)
	ErrUserAlreadyExists = topError.New("用户已存在", 1001)
	ErrUserNotFound      = topError.New("用户不存在", 1002)
	ErrUserRegFailed     = topError.New("用户注册失败", 1003)
	ErrLoginFailed       = topError.New("登录失败", 1004)
	ErrAuthFailed        = topError.New("认证失败", 1005)
	ErrInvalidPassword   = topError.New("密码错误", 1006)
	ErrUserNotActive     = topError.New("用户未激活", 1007)
	ErrNotDataSet        = topError.New("Not data set", 1005)

	// 模型相关错误
	ErrModelNotFound      = topError.New("模型不存在", 2001)
	ErrModelAlreadyExists = topError.New("模型已存在", 2002)
)
