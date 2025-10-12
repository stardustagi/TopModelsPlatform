package backend

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/stardustagi/TopLib/libs/jwt"
	"github.com/stardustagi/TopLib/libs/logs"
	"github.com/stardustagi/TopLib/libs/redis"
	"github.com/stardustagi/TopLib/utils"
	"github.com/stardustagi/TopModelsPlatform/constants"
	"go.uber.org/zap"
)

// PlatformUserAccess 节点用户访问
func PlatformUserAccess() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		//在这里处理拦截请求的逻辑
		return func(c echo.Context) error {
			//在这里处理拦截请求的逻辑
			jwtstr := c.Request().Header.Get("jwt")
			nodeUserId := c.Request().Header.Get("id")
			secret := fmt.Sprintf("%s-%s-%s", constants.AppName, constants.AppVersion, nodeUserId)
			if jwtstr != "" {
				jwtobj, ok := jwt.JWTDecrypt(jwtstr, secret)
				if !ok || jwtobj == nil || jwtobj["token"] == nil || jwtobj["id"] == nil {
					return c.JSON(401, map[string]interface{}{
						"errcode": 2,
						"errmsg":  "jwt解析错误",
					})
				}

				id, ok := jwtobj["id"].(string)
				if !ok {
					return c.JSON(401, map[string]interface{}{
						"errcode": 2,
						"errmsg":  "用户信息获取失败",
					})
				}
				intId, err := strconv.ParseInt(id, 10, 64)
				if err != nil {
					return c.JSON(401, map[string]interface{}{
						"errcode": 2,
						"errmsg":  "用户信息获取失败",
					})
				}
				tokenKey := fmt.Sprintf("%s:%s:user:%s", constants.AppName, constants.AppVersion, constants.PlatformUserTokenKey(intId))
				jwtKey := fmt.Sprintf("%s:%s", constants.ModelsKeyPrefix, constants.PlatformUserJwtKey(intId))
				logs.Info("redis key is:", zap.String("tokenKey", tokenKey), zap.String("jwtKey", jwtKey))
				redisCmd := redis.GetRedisDb()
				oldToken, err1 := redisCmd.Get(context.Background(), tokenKey).Result()
				if err1 != nil {
					return c.JSON(401, map[string]interface{}{
						"errcode": 2,
						"errmsg":  "获取token失败",
					})
				}
				if jwtobj["token"] != oldToken {
					return c.JSON(401, map[string]interface{}{
						"errcode": 2,
						"errmsg":  "账户已经在其他终端登录",
					})
				}
				//}
				c.QueryParams().Add("id", id)
			} else {
				ip := echo.ExtractIPFromXFFHeader()(c.Request())
				// 判断是否内网IP
				ok := utils.IsInnerIp(ip)
				if !ok {
					return c.JSON(401, map[string]interface{}{
						"errcode": 2,
						"errmsg":  "服务器繁忙,请稍后再试!",
					})
				}
			}
			return next(c)
		}
	}
}
