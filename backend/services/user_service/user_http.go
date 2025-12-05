package user_service

import (
	"context"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stardustagi/TopLib/libs/databases"
	"github.com/stardustagi/TopLib/libs/logs"
	"github.com/stardustagi/TopLib/libs/redis"
	"github.com/stardustagi/TopLib/libs/uuid"
	"github.com/stardustagi/TopLib/protocol"
	"github.com/stardustagi/TopModelsPlatform/backend"
	"github.com/stardustagi/TopModelsPlatform/constants"
	"github.com/stardustagi/TopModelsPlatform/models"
	"github.com/stardustagi/TopModelsPlatform/protocol/requests"
	"github.com/stardustagi/TopModelsPlatform/protocol/responses"
	"go.uber.org/zap"
)

type UserHttpService struct {
	logger    *zap.Logger
	ctx       context.Context
	cancelCtx context.CancelFunc
	dao       databases.BaseDao
	rds       redis.RedisCli
	app       *backend.Application
}

var (
	userHttpServiceInstance *UserHttpService
	userHttpServiceOnce     sync.Once
)

func GetUserHttpServiceInstance() *UserHttpService {
	userHttpServiceOnce.Do(func() {
		userHttpServiceInstance = NewUserHttpService()
	})
	return userHttpServiceInstance
}

func NewUserHttpService() *UserHttpService {
	ctx, cancel := context.WithCancel(context.Background())
	return &UserHttpService{
		logger:    logs.GetLogger("UserHttpService"),
		ctx:       ctx,
		cancelCtx: cancel,
		dao:       databases.GetDao(),
		rds: redis.NewRedisView(redis.GetRedisDb(),
			constants.ApplicationPrefix,
			logs.GetLogger("UserRedis")),
	}
}

func (u *UserHttpService) Start(app *backend.Application) {
	if app == nil {
		panic("请设置后端应用")
	}
	u.app = app
	u.initialization()
	u.logger.Info("Starting UserHttpService...")
}

func (u *UserHttpService) Stop() {
	u.logger.Info("Stopping UserHttpService...")
	u.cancelCtx()
	u.logger.Info("UserHttpService stopped.")
}

// List 获取用户列表
// @Summary 获取用户列表
// @Description 分页获取用户列表
// @Tags User
// @Accept json
// @Produce json
// @Param skip query int false "跳过条数"
// @Param limit query int false "每页条数"
// @Success 200 {object} responses.DefaultResponse
// @Router /user/list [get]
func (u *UserHttpService) List(c echo.Context, req requests.PageReq, resp responses.DefaultResponse) error {
	u.logger.Info("获取用户列表", zap.Int("skip", req.Skip), zap.Int("limit", req.Limit))

	session := u.dao.NewSession()
	defer session.Close()

	if req.Limit <= 0 {
		req.Limit = 20
	}
	if req.Limit > 100 {
		req.Limit = 100
	}

	var users []models.Users
	pageable := databases.NewPageable(req.Skip, req.Limit, req.Sort)
	total, err := session.FindAndCount(&users, pageable)
	if err != nil {
		u.logger.Error("查询用户列表失败", zap.Error(err))
		return protocol.Response(c, constants.ErrInternalServer, nil)
	}

	for i := range users {
		users[i].Password = ""
		users[i].Salt = ""
	}

	result := map[string]interface{}{
		"list":  users,
		"total": total,
		"skip":  req.Skip,
		"limit": req.Limit,
	}
	return protocol.Response(c, nil, result)
}

// GetById 根据ID获取用户
// @Summary 获取用户详情
// @Description 根据ID获取用户详情
// @Tags User
// @Accept json
// @Produce json
// @Param id query int true "用户ID"
// @Success 200 {object} responses.DefaultResponse
// @Router /user/get [get]
func (u *UserHttpService) GetById(c echo.Context, req requests.UserIdReq, resp responses.DefaultResponse) error {
	u.logger.Info("获取用户详情", zap.Int64("userId", req.Id))

	session := u.dao.NewSession()
	defer session.Close()

	user := &models.Users{Id: req.Id}
	ok, err := session.FindOne(user)
	if err != nil {
		u.logger.Error("查询用户失败", zap.Error(err))
		return protocol.Response(c, constants.ErrInternalServer, nil)
	}
	if !ok {
		return protocol.Response(c, constants.ErrUserNotFound, nil)
	}

	user.Password = ""
	user.Salt = ""
	return protocol.Response(c, nil, user)
}

// Create 创建用户
// @Summary 创建用户
// @Description 创建新用户
// @Tags User
// @Accept json
// @Produce json
// @Param request body requests.CreateUserReq true "创建用户请求"
// @Success 200 {object} responses.DefaultResponse
// @Router /user/create [post]
func (u *UserHttpService) Create(c echo.Context, req requests.CreateUserReq, resp responses.DefaultResponse) error {
	u.logger.Info("创建用户", zap.String("email", req.Email))

	session := u.dao.NewSession()
	defer session.Close()

	existUser := &models.Users{Email: req.Email}
	ok, _ := session.FindOne(existUser)
	if ok {
		return protocol.Response(c, constants.ErrUserAlreadyExists, nil)
	}

	salt := uuid.GenString(8)
	now := time.Now().Unix()

	user := &models.Users{
		UserName:  req.UserName,
		Email:     req.Email,
		Phone:     req.Phone,
		Password:  req.Password,
		Salt:      salt,
		Active:    1,
		CreatedAt: now,
	}

	_, err := session.InsertOne(user)
	if err != nil {
		u.logger.Error("创建用户失败", zap.Error(err))
		return protocol.Response(c, constants.ErrUserRegFailed, nil)
	}

	user.Password = ""
	user.Salt = ""
	return protocol.Response(c, nil, user)
}

// Update 更新用户
// @Summary 更新用户
// @Description 更新用户信息
// @Tags User
// @Accept json
// @Produce json
// @Param request body requests.UpdateUserReq true "更新用户请求"
// @Success 200 {object} responses.DefaultResponse
// @Router /user/update [post]
func (u *UserHttpService) Update(c echo.Context, req requests.UpdateUserReq, resp responses.DefaultResponse) error {
	u.logger.Info("更新用户", zap.Int64("userId", req.Id))

	session := u.dao.NewSession()
	defer session.Close()

	user := &models.Users{Id: req.Id}
	ok, err := session.FindOne(user)
	if err != nil || !ok {
		return protocol.Response(c, constants.ErrUserNotFound, nil)
	}

	if req.UserName != "" {
		user.UserName = req.UserName
	}
	if req.Phone != "" {
		user.Phone = req.Phone
	}
	if req.RealName != "" {
		user.RealName = req.RealName
	}
	user.LastUpdate = time.Now().Unix()

	_, err = session.UpdateById(user.Id, user)
	if err != nil {
		u.logger.Error("更新用户失败", zap.Error(err))
		return protocol.Response(c, constants.ErrInternalServer, nil)
	}

	user.Password = ""
	user.Salt = ""
	return protocol.Response(c, nil, user)
}

// Delete 删除用户
// @Summary 删除用户
// @Description 软删除用户
// @Tags User
// @Accept json
// @Produce json
// @Param request body requests.UserIdReq true "用户ID"
// @Success 200 {object} responses.DefaultResponse
// @Router /user/delete [post]
func (u *UserHttpService) Delete(c echo.Context, req requests.UserIdReq, resp responses.DefaultResponse) error {
	u.logger.Info("删除用户", zap.Int64("userId", req.Id))

	session := u.dao.NewSession()
	defer session.Close()

	user := &models.Users{Id: req.Id}
	ok, err := session.FindOne(user)
	if err != nil || !ok {
		return protocol.Response(c, constants.ErrUserNotFound, nil)
	}

	user.Deleted = 1
	user.LastUpdate = time.Now().Unix()

	_, err = session.UpdateById(user.Id, user)
	if err != nil {
		u.logger.Error("删除用户失败", zap.Error(err))
		return protocol.Response(c, constants.ErrInternalServer, nil)
	}

	return protocol.Response(c, nil, nil)
}

// SetActive 设置用户激活状态
// @Summary 设置用户激活状态
// @Description 设置用户激活或禁用
// @Tags User
// @Accept json
// @Produce json
// @Param request body requests.SetActiveReq true "设置激活状态请求"
// @Success 200 {object} responses.DefaultResponse
// @Router /user/setActive [post]
func (u *UserHttpService) SetActive(c echo.Context, req requests.SetActiveReq, resp responses.DefaultResponse) error {
	u.logger.Info("设置用户激活状态", zap.Int64("userId", req.Id), zap.Int("active", req.Active))

	session := u.dao.NewSession()
	defer session.Close()

	user := &models.Users{Id: req.Id}
	ok, err := session.FindOne(user)
	if err != nil || !ok {
		return protocol.Response(c, constants.ErrUserNotFound, nil)
	}

	user.Active = req.Active
	user.LastUpdate = time.Now().Unix()

	_, err = session.UpdateById(user.Id, user)
	if err != nil {
		u.logger.Error("设置用户激活状态失败", zap.Error(err))
		return protocol.Response(c, constants.ErrInternalServer, nil)
	}

	return protocol.Response(c, nil, nil)
}
