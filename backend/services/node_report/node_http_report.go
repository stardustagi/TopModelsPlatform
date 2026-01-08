package node_report

import (
	"context"
	"sync"

	"github.com/stardustagi/TopLib/libs/databases"
	"github.com/stardustagi/TopLib/libs/logs"
	"github.com/stardustagi/TopLib/libs/redis"
	"github.com/stardustagi/TopLib/libs/server"
	"github.com/stardustagi/TopModelsPlatform/backend"
	"github.com/stardustagi/TopModelsPlatform/constants"
	"go.uber.org/zap"
)

type NodeHttpReportService struct {
	logger    *zap.Logger
	ctx       context.Context
	cancelCtx context.CancelFunc
	dao       databases.BaseDao
	rds       redis.RedisCli
	mu        sync.RWMutex
	app       *backend.Application
}

var (
	nodeHttpReportServiceInstance *NodeHttpReportService
	nodeHttpReportServiceOnce     sync.Once
)

// GetNodeHttpReportServiceInstance 获取报表服务实例
func GetNodeHttpReportServiceInstance() *NodeHttpReportService {
	nodeHttpReportServiceOnce.Do(func() {
		nodeHttpReportServiceInstance = NewNodeHttpReportService()
	})
	return nodeHttpReportServiceInstance
}

// NewNodeHttpReportService 创建新的报表服务
func NewNodeHttpReportService() *NodeHttpReportService {
	ctx, cancel := context.WithCancel(context.Background())
	return &NodeHttpReportService{
		logger:    logs.GetLogger("NodeHttpReportService"),
		ctx:       ctx,
		cancelCtx: cancel,
		dao:       databases.GetDao(),
		rds: redis.NewRedisView(redis.GetRedisDb(),
			constants.PlatformPrefix,
			logs.GetLogger("NodeReportRedis")),
	}
}

func (rep *NodeHttpReportService) Start(app *backend.Application) {
	if app == nil {
		panic("请设置后端应用")
	}
	rep.app = app
	rep.initialization()
	rep.logger.Info("Starting NodeHttpReportService...")

}

func (rep *NodeHttpReportService) Stop() {
	rep.logger.Info("Stopping NodeHttpReportService...")
	rep.cancelCtx()
	rep.logger.Info("NodeHttpReportService stopped.")
}
