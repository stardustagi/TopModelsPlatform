package platform

import (
	"context"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stardustagi/TopLib/libs/databases"
	"github.com/stardustagi/TopLib/libs/logs"
	"github.com/stardustagi/TopLib/libs/redis"
	"github.com/stardustagi/TopLib/protocol"
	"github.com/stardustagi/TopModelsPlatform/backend"
	"github.com/stardustagi/TopModelsPlatform/constants"
	"github.com/stardustagi/TopModelsPlatform/models"
	"github.com/stardustagi/TopModelsPlatform/protocol/requests"
	"github.com/stardustagi/TopModelsPlatform/protocol/responses"

	"go.uber.org/zap"
)

type PlatfromHttpService struct {
	logger      *zap.Logger
	ctx         context.Context
	cancelCtx   context.CancelFunc
	dao         databases.BaseDao
	rds         redis.RedisCli
	mu          sync.RWMutex // 读写锁保护usersJwt
	userInfosMu sync.RWMutex // 读写保护UserInfos
	app         *backend.Application
}

var (
	platformHttpServiceInstance *PlatfromHttpService
	platofrmHttpServiceOnce     sync.Once
)

// GetUsersHttpServiceInstance 获取 HTTP 用户服务实例
func GetUsersHttpServiceInstance() *PlatfromHttpService {
	platofrmHttpServiceOnce.Do(func() {
		platformHttpServiceInstance = NewPlatformHttpService()
	})
	return platformHttpServiceInstance
}

// NewPlatformHttpService 创建新的 HTTP 用户服务
func NewPlatformHttpService() *PlatfromHttpService {
	ctx, cancel := context.WithCancel(context.Background())
	return &PlatfromHttpService{
		logger:    logs.GetLogger("UsersHttpService"),
		ctx:       ctx,
		cancelCtx: cancel,
		dao:       databases.GetDao(),
		rds: redis.NewRedisView(redis.GetRedisDb(),
			constants.ApplicationPrefix,
			logs.GetLogger("NodeUserRedis")),
	}
}

func (p *PlatfromHttpService) Start(app *backend.Application) {
	if app == nil {
		panic("请设置后端应用")
	}
	p.app = app
	p.initialization()
	p.logger.Info("Starting UsersHttpService...")
}

func (p *PlatfromHttpService) Stop() {
	p.logger.Info("Stopping UsersHttpService...")
	p.cancelCtx()
	p.logger.Info("UsersHttpService stopped.")
}

func (p *PlatfromHttpService) Register(ctx echo.Context,
	req requests.RegisterReq, resp requests.DefaultRequest) error {
	return nil
}

// Login 用户登录
// @Summary 用户登录
// @Description 用户登录接口
// @Tags Platform
// @Accept json
// @Produce json
// @Param request body requests.LoginReq true "登录请求"
// @Success 200 {object} responses.DefaultResponse
// @Router /platform/login [post]
func (p *PlatfromHttpService) Login(c echo.Context, req requests.LoginReq, resp responses.DefaultResponse) error {
	p.logger.Info("Platform user login request", zap.String("email", req.Mail))

	// 查询用户是否存在
	session := p.dao.NewSession()
	defer session.Close()

	platformUser := &models.PlatformUsers{Email: req.Mail}
	ok, err := session.FindOne(platformUser)
	if err != nil {
		p.logger.Error("查询平台用户失败", zap.Error(err), zap.String("email", req.Mail))
		return protocol.Response(c, constants.ErrInternalServer, nil)
	}
	if !ok {
		p.logger.Warn("平台用户不存在", zap.String("email", req.Mail))
		return protocol.Response(c, constants.ErrUserNotFound, nil)
	}

	// 检查用户是否已激活
	if platformUser.Active != 1 {
		p.logger.Warn("平台用户未激活", zap.String("email", req.Mail))
		return protocol.Response(c, constants.ErrUserNotActive, nil)
	}

	// 验证密码 - 使用 nodeUserMailDecodeToken 解密验证
	_, err = p.nodeUserMailDecodeToken(req.Password, platformUser.Password, platformUser.Salt)
	if err != nil {
		p.logger.Error("密码验证失败", zap.Error(err), zap.String("email", req.Mail))
		return protocol.Response(c, constants.ErrInvalidPassword, nil)
	}

	// 生成 JWT token
	token, err := p.generateUserToken(platformUser.Id, platformUser.Email)
	if err != nil {
		p.logger.Error("生成token失败", zap.Error(err), zap.String("email", req.Mail))
		return protocol.Response(c, constants.ErrInternalServer, nil)
	}

	// 更新最后登录时间
	platformUser.LastUpdate = time.Now().Unix()
	_, err = session.UpdateById(platformUser.Id, platformUser)
	if err != nil {
		p.logger.Warn("更新用户最后登录时间失败", zap.Error(err), zap.Int64("userId", platformUser.Id))
	}

	// 将token存储到Redis
	key := constants.PlatformUserTokenKey(platformUser.Id)
	err = p.rds.Set(p.ctx, key, []byte(token), "24h")
	if err != nil {
		p.logger.Warn("保存token到Redis失败", zap.Error(err), zap.Int64("userId", platformUser.Id))
	}

	p.logger.Info("平台用户登录成功", zap.String("email", req.Mail), zap.Int64("userId", platformUser.Id))

	// 返回登录成功信息
	result := map[string]interface{}{
		"token":    token,
		"userId":   platformUser.Id,
		"email":    platformUser.Email,
		"userName": platformUser.UserName,
	}

	return protocol.Response(c, nil, result)
}

// RegisterUser 用户注册
// @Summary 用户注册
// @Description 平台用户注册接口
// @Tags Platform
// @Accept json
// @Produce json
// @Param request body requests.RegisterReq true "注册请求"
// @Success 200 {object} responses.DefaultResponse
// @Router /platform/register [post]
func (p *PlatfromHttpService) RegisterUser(c echo.Context, req requests.RegisterReq, resp responses.DefaultResponse) error {
	p.logger.Info("Platform user register request", zap.String("email", req.Mail))

	session := p.dao.NewSession()
	defer session.Close()

	// 检查用户是否已存在
	existUser := &models.PlatformUsers{Email: req.Mail}
	ok, _ := session.FindOne(existUser)
	if ok {
		return protocol.Response(c, constants.ErrUserAlreadyExists, nil)
	}

	// 生成盐值
	salt := "salt1234" // 实际应使用 uuid.GenString(8)
	now := time.Now().Unix()

	// 加密密码
	encryptedPwd, err := p.nodeUserMailEncodeToken(req.Mail, req.Password, salt)
	if err != nil {
		p.logger.Error("密码加密失败", zap.Error(err))
		return protocol.Response(c, constants.ErrInternalServer, nil)
	}

	user := &models.PlatformUsers{
		Email:      req.Mail,
		Password:   encryptedPwd,
		Salt:       salt,
		Active:     1,
		CreatedAt:  now,
		LastUpdate: now,
	}

	_, err = session.InsertOne(user)
	if err != nil {
		p.logger.Error("创建用户失败", zap.Error(err))
		return protocol.Response(c, constants.ErrUserRegFailed, nil)
	}

	p.logger.Info("平台用户注册成功", zap.String("email", req.Mail), zap.Int64("userId", user.Id))

	result := map[string]interface{}{
		"userId": user.Id,
		"email":  user.Email,
	}
	return protocol.Response(c, nil, result)
}
