package node_report

import (
	"context"
	"sync"

	"github.com/stardustagi/TopLib/libs/databases"
	"github.com/stardustagi/TopLib/libs/logs"
	"github.com/stardustagi/TopLib/libs/redis"
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

	isRunning bool
	stopCh    chan struct{}
	wg        sync.WaitGroup
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
		stopCh:    make(chan struct{}),
		rds: redis.NewRedisView(redis.GetRedisDb(),
			constants.PlatformPrefix,
			logs.GetLogger("NodeReportRedis")),
	}
}

func (n *NodeHttpReportService) Start(app *backend.Application) {
	if app == nil {
		panic("请设置后端应用")
	}
	n.app = app
	n.logger.Info("Starting NodeHttpReportService...")

	n.mu.Lock()
	defer n.mu.Unlock()
	if n.isRunning {
		n.logger.Info("NodeHttpReportService already running")
		return
	}

	n.isRunning = true
	n.logger.Info("NodeHttpReportService started")
}

func (n *NodeHttpReportService) Stop() {
	n.logger.Info("Stopping NodeHttpReportService...")
	n.cancelCtx()

	n.mu.Lock()
	defer n.mu.Unlock()
	if !n.isRunning {
		n.logger.Info("NodeHttpReportService already stopped")
		return
	}

	close(n.stopCh)
	n.wg.Wait()

	n.stopCh = make(chan struct{})
	n.isRunning = false
	n.logger.Info("NodeHttpReportService stopped.")
}
