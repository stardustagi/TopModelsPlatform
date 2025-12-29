package platform

import (
	"context"
	"fmt"
	"time"

	"github.com/stardustagi/TopLib/libs/jwt"
	"github.com/stardustagi/TopLib/libs/redis"
	"github.com/stardustagi/TopLib/libs/uuid"
	"github.com/stardustagi/TopModelsPlatform/constants"
	"go.uber.org/zap"
)

// nodeUserMailEncodeToken 节点用户邮箱加密
func (p *PlatfromHttpService) nodeUserMailEncodeToken(mail, password string, fixedSalt string) (string, error) {
	return jwt.EncryptWithFixedSalt(password, constants.ObtenationIterations, mail, fixedSalt)
}

// nodeUserMailDecodeToken 节点用户邮箱解密
func (p *PlatfromHttpService) nodeUserMailDecodeToken(inPassword, dbPassword string, fixedSalt string) (string, error) {
	return jwt.DecryptWithFixedSalt(inPassword, constants.ObtenationIterations, dbPassword, fixedSalt)
}

// generateUserToken 生成用户Token (简化版)
func (p *PlatfromHttpService) generateUserToken(userId int64, email string) (string, error) {
	return p.generatePlatformUserToken(userId, fmt.Sprintf("%d", userId), email)
}

// generatePlatformUserToken 生成平台用户Token
// 参照 generateNodeUserToken 实现，与 app_middleware.go 中的验证逻辑对齐
func (p *PlatfromHttpService) generatePlatformUserToken(userId int64, once, email string) (string, error) {
	p.logger.Info("生成平台用户Token", zap.Int64("userId", userId), zap.String("email", email))

	// 使用原生Redis客户端，与 app_middleware.go 保持一致
	redisCmd := redis.GetRedisDb()
	ctx := context.Background()

	// 生成用户Token Key - 与 app_middleware.go 中的 tokenKey 格式完全对齐
	// middleware: tokenKey := fmt.Sprintf("%s:%s:user:%s", constants.AppName, constants.AppVersion, constants.PlatformUserTokenKey(intId))
	userTokenKey := fmt.Sprintf("%s:%s:user:%s", constants.AppName, constants.AppVersion, constants.PlatformUserTokenKey(userId))

	// 生成随机token
	token := uuid.GenString(32)

	// 解析过期时间
	expiration, err := time.ParseDuration(constants.PlatformUserTokenExpire)
	if err != nil {
		p.logger.Error("解析Token过期时间失败", zap.Error(err))
		expiration = 720 * time.Hour // 默认720小时
	}

	// 存储token到Redis - 使用原生客户端，不带view前缀
	err = redisCmd.Set(ctx, userTokenKey, token, expiration).Err()
	if err != nil {
		p.logger.Error("保存平台用户Token到Redis失败", zap.Error(err), zap.String("userTokenKey", userTokenKey))
		return "", err
	}

	// 组装JWT密钥 - 与 app_middleware.go 中的 secret 格式对齐
	// middleware: secret := fmt.Sprintf("%s-%s-%s", constants.AppName, constants.AppVersion, UserId)
	userIdStr := fmt.Sprintf("%d", userId)
	jwtKey := fmt.Sprintf("%s-%s-%s", constants.AppName, constants.AppVersion, once)

	// 生成JWT Token
	jwtToken := jwt.JWTEncrypt(userIdStr, token, jwtKey)

	// 存储JWT到Redis - 与 app_middleware.go 中的 jwtKey 格式对齐
	// middleware: jwtKey := fmt.Sprintf("%s:%s", constants.ModelsKeyPrefix, constants.PlatformUserJwtKey(intId))
	jwtTokenKey := fmt.Sprintf("%s:%s", constants.ModelsKeyPrefix, constants.PlatformUserJwtKey(userId))
	err = redisCmd.Set(ctx, jwtTokenKey, jwtToken, expiration).Err()
	if err != nil {
		p.logger.Error("保存平台用户JWT到Redis失败", zap.Error(err), zap.String("jwtTokenKey", jwtTokenKey))
		// 清理之前保存的token
		redisCmd.Del(ctx, userTokenKey)
		return "", err
	}

	p.logger.Info("平台用户Token生成成功",
		zap.Int64("userId", userId),
		zap.String("userTokenKey", userTokenKey),
		zap.String("jwtTokenKey", jwtTokenKey))
	return jwtToken, nil
}
