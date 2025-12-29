package constants

import (
	"fmt"
	"os"
)

var (
	ApplicationName   = "TopModelsPlatform"
	ApplicationPrefix = "platform"
)

var (
	WechatOpenIDKey      string = "wechat:openid"
	UserDefaultPassword  string = "FsNyBfcpdtq1p0063MBu"
	PhoneDefaultPassword string = "6R0b8zfcl6rBb5vkeqVj"
	MailDefaultPassword  string = "ajZxPKxauQaHOtpBJF6W"
	ObtenationIterations int    = 3
	CodeKeyExpire        int    = 30
	UserGraphCodeExpire  string = "300" // 图形验证码过期时间（秒）
	UserPhoneCodeExpire  int    = 300   // 手机验证码过期时间（秒）
	PlatformPrefix       string = fmt.Sprintf("%s:platform", RedisPrefix)
)

var (
	HeaderId  string = "id"
	JwtSecret string
)

var (
	AppName    string
	AppVersion string
)

func Init() {
	AppName = os.Getenv("APP_NAME")
	AppVersion = os.Getenv("APP_VERSION")
	RedisPrefix = fmt.Sprintf("%s:%s", AppName, AppVersion)
	ModelsKeyPrefix = fmt.Sprintf("%s:llm:models", RedisPrefix)
}

// PlatformUserTokenKey 节点用户TokenKey
func PlatformUserTokenKey(id int64) string {
	return fmt.Sprintf("nodeUserToken:%d", id)
}

func PlatformUserJwtKey(id int64) string {
	return fmt.Sprintf("nodeUserJwt:%d", id)
}

func PlatformUserGraphVerifyKey(t string) string {
	return fmt.Sprintf("nodeUserGraphVerifyCode:%s", t)
}

func PlatformUserPhoneVerifyKey(phone string) string {
	return fmt.Sprintf("nodeUserPhoneVerifyCode:%s", phone)
}

func PlatformUserEmailVerifyKey(email string) string {
	return fmt.Sprintf("nodeUserEmailVerifyCode:%s", email)
}

func PlatformAccessKey(nodeUserId int64, ak string) string {
	return fmt.Sprintf("%d:%s", nodeUserId, ak)
}

func PlatformUserMailVerifyKey(nodeUserId int64) string {
	return fmt.Sprintf("mail:verify:%d", nodeUserId)
}

func PlatformAccessModelsKey(nodeId string) string {
	return fmt.Sprintf("node:access:%s", nodeId)
}

func PlatformUserAccessTokenKey(nodeUserId int64) string {
	return fmt.Sprintf("nodeUser:token:%d", nodeUserId)
}
