package platform

import (
	"fmt"

	"github.com/stardustagi/TopLib/libs/jwt"
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
// 参照 generateNodeUserToken 实现
func (p *PlatfromHttpService) generatePlatformUserToken(userId int64, once, email string) (string, error) {
	p.logger.Info("生成平台用户Token", zap.Int64("userId", userId), zap.String("email", email))

	// 生成用户Token Key
	userTokenKey := constants.PlatformUserTokenKey(userId)

	// 生成随机token
	token := uuid.GenBytes(32)

	// 存储token到Redis
	err := p.rds.Set(p.ctx, userTokenKey, token, constants.PlatformUserTokenExpire)
	if err != nil {
		p.logger.Error("保存平台用户Token到Redis失败", zap.Error(err), zap.String("userTokenKey", userTokenKey))
		return "", err
	}

	// 组装JWT密钥
	userIdStr := fmt.Sprintf("%d", userId)
	jwtKey := fmt.Sprintf("%s-%s-%s", constants.AppName, constants.AppVersion, once)

	// 生成JWT Token
	jwtToken := jwt.JWTEncrypt(userIdStr, string(token), jwtKey)

	// 存储JWT到Redis
	jwtTokenKey := constants.PlatformUserJwtKey(userId)
	err = p.rds.Set(p.ctx, jwtTokenKey, []byte(jwtToken), constants.PlatformUserTokenExpire)
	if err != nil {
		p.logger.Error("保存平台用户JWT到Redis失败", zap.Error(err), zap.String("jwtTokenKey", jwtTokenKey))
		// 清理之前保存的token
		_, _ = p.rds.Del(p.ctx, userTokenKey)
		return "", err
	}

	p.logger.Info("平台用户Token生成成功", zap.Int64("userId", userId))
	return jwtToken, nil
}
