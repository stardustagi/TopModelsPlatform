package platform

import (
	"errors"
	"fmt"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stardustagi/TopLib/libs/redis"
	"github.com/stardustagi/TopLib/protocol"
	"github.com/stardustagi/TopModelsPlatform/constants"
	"github.com/stardustagi/TopModelsPlatform/protocol/requests"
	"github.com/stardustagi/TopModelsPlatform/protocol/responses"
	"go.uber.org/zap"
)

// GetGraphCode 获取图形验证码
// @Summary 获取图形验证码
// @Description 通过类型获取图形验证码
// @Tags System
// @Accept json
// @Produce json
// @Param request body requests.GetGraphVerifyCodeReq true "获取图形验证码请求"
// @Success 200 {object} responses.GetGraphVerifyCodeResp "成功"
// @Failure 400 {object} responses.DefaultResponse "参数错误"
// @Failure 500 {object} responses.DefaultResponse "服务器错误"
// @Router /system/graphCode [post]
func (p *PlatfromHttpService) GetGraphCode(c echo.Context,
	req requests.GetGraphVerifyCodeReq, resp responses.GetGraphVerifyCodeResp) error {
	_t := req.T
	if _t == "" {
		return protocol.Response(c, constants.ErrVerifyCode, nil)
	}

	code, err := p.GetGraphVerifyCodeKey(_t)
	if err != nil {
		p.logger.Error("Failed to get graph code", zap.Error(err), zap.String("type", _t))
		return protocol.Response(c, constants.ErrGetVerifyCode, nil)
	}
	p.logger.Info("Verify code generated", zap.String("key", code))

	resp.Code = code
	resp.ExpireTime = constants.UserGraphCodeExpire
	return protocol.Response(c, nil, resp)
}

// GetPhoneVerifyCode 获取手机验证码
// @Summary 获取手机验证码
// @Description 通过手机号和图形验证码获取手机验证码
// @Tags System
// @Accept json
// @Produce json
// @Param request body requests.GetPhoneVerifyCodeReq true "获取手机验证码请求"
// @Success 200 {object} responses.DefaultResponse "成功"
// @Failure 400 {object} responses.DefaultResponse "参数错误"
// @Failure 500 {object} responses.DefaultResponse "服务器错误"
// @Router /system/phoneCode [post]
func (p *PlatfromHttpService) GetPhoneCode(c echo.Context,
	req requests.GetPhoneVerifyCodeReq, resp responses.DefaultResponse) error {
	_t := req.T
	if _t == "" {
		return protocol.Response(c, constants.ErrInvalidParams, nil)
	}

	// 验证图形码,生成手机验证码
	code, err := p.generatePhoneVerifyCode(req.Phone, req.GraphCode, _t)
	if err != nil && !errors.Is(err, redis.Nil) {
		p.logger.Error("Failed to verify graph code", zap.Error(err), zap.String("type", _t))
		return protocol.Response(c, constants.ErrVerifyCode, nil)
	}

	phone := req.Phone
	if phone == "" {
		p.logger.Error("Phone number is required")
		return protocol.Response(c, constants.ErrInvalidParams, nil)
	}

	// TODO: 调用短信网关发送验证码
	// msgBody := fmt.Sprintf("{\"code\": \"%s\"}", code)
	// _, err = sms_gateway.SendByteSMS(phone, msgBody)
	// if err != nil {
	// 	s.logger.Error("Failed to send SMS", zap.Error(err), zap.String("phone", phone))
	// 	return protocol.Response(c, constants.ErrSendSMS.AppendErrors(err), nil)
	// }

	p.logger.Info("verify code generated", zap.String("type", req.RegType), zap.String("key", code))
	return protocol.Response(c, nil, fmt.Sprintf("验证码:%s 已发送，请注意查收, 5分钟内有效", code))
}

// AdminLogin 管理员用户登录
// @Summary 管理员用户登录
// @Description 管理员用户登录
// @Tags System
// @Accept json
// @Produce json
// @Param data body requests.PhoneLoginAndRegisterReq true "登录请求数据"
// @Success 200 {object} responses.UserLoginAndRegisterRes "成功返回数据"
// @Failure 400 {object} responses.DefaultResponse "请求参数错误"
// @Failure 500 {object} responses.DefaultResponse "服务器内部错误"
// @Router /system/adminLogin [post]
//func (p *PlatfromHttpService) AdminLogin(c echo.Context,
//	req requests.PhoneLoginAndRegisterReq, resp responses.UserLoginAndRegisterRes) error {
//	p.logger.Info("Admin login requested", zap.String("phone", req.Phone))
//
//	clientIP := c.RealIP()
//	if clientIP == "" {
//		p.logger.Error("Client IP is required")
//		return protocol.Response(c, constants.ErrInvalidParams, nil)
//	}
//
//	// TODO: 实现管理员登录逻辑
//	// result, err := s.AdminPhoneLogin(req.Phone, req.PhoneCode, req.T, req.InviteCode, clientIP)
//	// if err != nil {
//	// 	s.logger.Error("Failed to admin login", zap.Error(err))
//	// 	return protocol.Response(c, constants.ErrInternalServer.AppendErrors(err), nil)
//	// }
//
//	// 设置响应头
//	// c.Response().Header().Set("Access-Control-Expose-Headers", "id, jwt")
//	// c.Response().Header().Set("id", fmt.Sprintf("%d", result.UserID))
//	// c.Response().Header().Set("jwt", result.Token)
//
//	return protocol.Response(c, nil, resp)
//}

// GetGraphVerifyCodeKey 生成图形验证码并存储到Redis
func (p *PlatfromHttpService) GetGraphVerifyCodeKey(t string) (string, error) {
	// 生成6位数字验证码
	code := fmt.Sprintf("%06d", p.generateRandomCode())
	key := constants.PlatformUserGraphVerifyKey(t)

	err := p.rds.Set(p.ctx, key, []byte(code), fmt.Sprintf("%d", constants.UserPhoneCodeExpire))
	if err != nil {
		return "", err
	}
	return code, nil
}

// generatePhoneVerifyCode 验证图形验证码并生成手机验证码
func (p *PlatfromHttpService) generatePhoneVerifyCode(phone, graphCode, t string) (string, error) {
	// 验证图形验证码
	key := constants.PlatformUserGraphVerifyKey(t)
	storedCode, err := p.rds.Get(p.ctx, key)
	if err != nil {
		return "", err
	}
	if string(storedCode) != graphCode {
		return "", fmt.Errorf("图形验证码错误")
	}

	// 生成手机验证码
	code := fmt.Sprintf("%06d", p.generateRandomCode())
	phoneKey := constants.PlatformUserPhoneVerifyKey(phone)

	err = p.rds.Set(p.ctx, phoneKey, []byte(code), fmt.Sprintf("%d", constants.UserPhoneCodeExpire))
	if err != nil {
		return "", err
	}
	return code, nil
}

func (p *PlatfromHttpService) generateRandomCode() int {
	// 生成100000-999999之间的随机数
	return 100000 + int(time.Now().UnixNano()%900000)
}
