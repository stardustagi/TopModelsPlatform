package main

import (
	"github.com/stardustagi/TopLib/libs/conf"
	"github.com/stardustagi/TopLib/libs/databases"
	"github.com/stardustagi/TopLib/libs/logs"
	"github.com/stardustagi/TopLib/libs/redis"
	"github.com/stardustagi/TopModelsPlatform/backend"
	platform "github.com/stardustagi/TopModelsPlatform/backend/services/platform_service"
	user_service "github.com/stardustagi/TopModelsPlatform/backend/services/user_service"
	"github.com/stardustagi/TopModelsPlatform/constants"

	_ "github.com/stardustagi/TopModelsPlatform/docs"
)

// @title TopModelsPlatform API
// @version 1.0
// @description LLM 管理平台 API 文档
// @host localhost:8080
// @BasePath /
func main() {
	conf.Init()
	loggerConfig := conf.Get("logger")
	constants.Init()
	logs.Init(loggerConfig)
	logger := logs.GetLogger("main")
	logger.Info("Init logs")
	_, _ = databases.Init(conf.Get("mysql"))
	logger.Info("Init mysql")
	_, _ = redis.Init(conf.Get("redis"))
	logger.Info("Init redis")

	app := backend.NewApplication(conf.Get("websrv"), conf.Get("websocket"))

	// 启动平台服务
	platformService := platform.GetUsersHttpServiceInstance()
	platformService.Start(app)
	logger.Info("Platform service started")

	// 启动用户服务
	userService := user_service.GetUserHttpServiceInstance()
	userService.Start(app)
	logger.Info("User service started")

	app.Start()

	// 停止服务
	platformService.Stop()
	userService.Stop()
	app.Stop()
}
