package backend

import (
	"context"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/stardustagi/TopLib/libs/databases"
	"github.com/stardustagi/TopLib/libs/logs"
	"github.com/stardustagi/TopLib/libs/server"
	"github.com/stardustagi/TopLib/utils"
	"github.com/stardustagi/TopModelsPlatform/models"
	echoSwagger "github.com/swaggo/echo-swagger"
	"go.uber.org/zap"
)

type Application struct {
	ctx      context.Context
	b        *server.Backend
	logger   *zap.Logger
	config   server.HttpServerConfig
	wsConfig server.HttpWebSocketConfig
	manager  server.IClientManager
	upgrader websocket.Upgrader
}

func NewApplication(configBytes, wsConfigBytes []byte) *Application {
	config, err := utils.Bytes2Struct[server.HttpServerConfig](configBytes)
	if err != nil {
		panic(err)
	}
	wsConfig, err := utils.Bytes2Struct[server.HttpWebSocketConfig](wsConfigBytes)
	if err != nil {
		panic(err)
	}
	b, err := server.NewBackend(configBytes)
	app := &Application{
		ctx:      context.Background(),
		config:   config,
		wsConfig: wsConfig,
		logger:   logs.GetLogger("HttpBackend"),
		manager:  server.NewClientManager(logs.GetLogger("clientManager")),
		upgrader: websocket.Upgrader{
			ReadBufferSize:  wsConfig.ReadBufferSize,
			WriteBufferSize: wsConfig.WriteBufferSize,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		b: b,
	}
	// 注册 swagger 路由
	app.b.AddNativeHandler("GET", "/swagger/*", echoSwagger.WrapHandler)
	return app
}

func (h *Application) Start() {
	h.syncDatabaseSchema()
	go h.manager.Start()
	go func() {
		if err := h.b.Start(); err != nil {
			h.logger.Error("backend.Start error", zap.Error(err))
		}
	}()
}

func (h *Application) Stop() {
	h.logger.Info("Stopping HttpBackend")
	h.b.Stop()
	h.manager.Stop()
}

func (h *Application) AddGroup(group string, middleware ...echo.MiddlewareFunc) {
	h.b.AddGroup(group, middleware...)
}

func (h *Application) AddPostHandler(group string, handler server.IHandler) {
	h.b.AddPostHandler(group, handler)
}

func (h *Application) AddGetHandler(group string, handler server.IHandler) {
	h.b.AddGetHandler(group, handler)
}

func (h *Application) AddNativeHandler(method, path string, handler echo.HandlerFunc) {
	h.b.AddNativeHandler(method, path, handler)
}

func (h *Application) HandleWebSocket() {
}

func (h *Application) syncDatabaseSchema() {
	h.logger.Info("Syncing database schema...")
	modelList := []interface{}{
		&models.PlatformUsers{},
		&models.ProviderConsumeSummary{},
		&models.ProviderModelDailySummary{},
		&models.UserRebateConfig{},
		&models.UserRebateMonthly{},
		&models.UserWallet{},
	}
	dbDao := databases.GetDao()
	if err := dbDao.Native().Sync2(modelList...); err != nil {
		h.logger.Error("Database schema sync failed", zap.Error(err))
		panic("Database schema sync failed: " + err.Error())
	}
	h.logger.Info("Database schema synced successfully")
}
