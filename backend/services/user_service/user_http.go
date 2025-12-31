package user_service

import (
	"context"
	"fmt"
	"strconv"
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
	if req.Sort == "" {
		req.Sort = "id ASC"
	}

	var users []models.Users

	// 使用 Native().Where() 显式过滤已删除用户，避免零值被忽略
	total, err := session.Native().
		Where("deleted = ?", 0).
		OrderBy(req.Sort).
		Limit(req.Limit, req.Skip).
		FindAndCount(&users)
	if err != nil {
		u.logger.Error("查询用户列表失败", zap.Error(err))
		return protocol.Response(c, constants.ErrInternalServer, nil)
	}

	// 构建返回结果，包含钱包余额信息
	type UserWithWallet struct {
		models.Users
		Balance       int64 `json:"balance"`
		RebateBalance int64 `json:"rebate_balance"`
	}

	var userList []UserWithWallet
	for _, user := range users {
		user.Password = ""
		user.Salt = ""

		// 查询用户钱包信息
		wallet := &models.UserWallet{UserId: user.Id}
		ok, _ := session.FindOne(wallet)

		userWithWallet := UserWithWallet{
			Users:         user,
			Balance:       0,
			RebateBalance: 0,
		}
		if ok {
			userWithWallet.Balance = wallet.Balance
			userWithWallet.RebateBalance = wallet.RebateBalance
		}
		userList = append(userList, userWithWallet)
	}

	result := map[string]interface{}{
		"list":  userList,
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

// AdminAddUserByName 管理员添加用户
// @Summary 管理员添加用户
// @Description 管理员使用用户名添加用户
// @Tags User
// @Accept json
// @Produce json
// @Param data body requests.AdminAddUserReq true "管理员添加用户请求数据"
// @Success 200 {object} responses.DefaultResponse "成功返回数据"
// @Failure 400 {object} responses.DefaultResponse "请求参数错误"
// @Failure 500 {object} responses.DefaultResponse "服务器内部错误"
// @Router /user/adminAddUserByName [post]
func (u *UserHttpService) AdminAddUserByName(c echo.Context, req requests.AdminAddUserReq, resp responses.DefaultResponse) error {
	u.logger.Info("管理员添加用户", zap.String("userName", req.UserName), zap.String("phone", req.Phone))

	session := u.dao.NewSession()
	defer session.Close()

	// 检查用户名是否存在
	existUser := &models.Users{UserName: req.UserName}
	ok, _ := session.FindOne(existUser)
	if ok {
		return protocol.Response(c, constants.ErrUserAlreadyExists, nil)
	}

	// 检查手机号是否存在
	existPhone := &models.Users{Phone: req.Phone}
	ok, _ = session.FindOne(existPhone)
	if ok {
		return protocol.Response(c, constants.ErrUserAlreadyExists.AppendErrors(fmt.Errorf("手机号已存在")), nil)
	}

	salt := uuid.GenString(8)
	now := time.Now().Unix()

	user := &models.Users{
		UserName:  req.UserName,
		Phone:     req.Phone,
		Password:  req.Password, // TODO: 应该加密存储
		Salt:      salt,
		Active:    1,
		CreatedAt: now,
	}

	_, err := session.InsertOne(user)
	if err != nil {
		u.logger.Error("管理员添加用户失败", zap.Error(err))
		return protocol.Response(c, constants.ErrUserRegFailed.AppendErrors(err), nil)
	}

	return protocol.Response(c, nil, map[string]interface{}{
		"success": true,
		"message": "使用用户名添加用户成功",
		"user_id": user.Id,
	})
}

// UserAdminPayment 用户充值接口
// @Summary 用户充值接口
// @Description 管理员为用户充值
// @Tags User
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer {token}"
// @Param req body requests.UserAdminPaymentReq true "请求参数"
// @Success 200 {object} responses.DefaultResponse
// @Failure 400 {object} responses.DefaultResponse
// @Failure 500 {object} responses.DefaultResponse
// @Router /user/userAdminPayment [post]
func (u *UserHttpService) UserAdminPayment(c echo.Context, req requests.UserAdminPaymentReq, resp responses.DefaultResponse) error {
	u.logger.Info("用户充值", zap.String("userName", req.UserName), zap.Int64("amount", req.Amount))

	adminUserIdStr := c.Request().Header.Get("id")
	adminUserId, err := strconv.ParseInt(adminUserIdStr, 10, 64)
	if err != nil {
		return protocol.Response(c, constants.ErrInvalidParams.AppendErrors(err), nil)
	}

	// 检查管理员权限
	ok, err := u.checkUserIsAdmin(adminUserId)
	if err != nil {
		return protocol.Response(c, constants.ErrInternalServer.AppendErrors(err), nil)
	}
	if !ok {
		u.logger.Info("not admin user cant do this", zap.Int64("userId", adminUserId))
		return protocol.Response(c, constants.ErrAuthFailed, nil)
	}

	session := u.dao.NewSession()
	defer session.Close()

	// 查找用户
	user := &models.Users{UserName: req.UserName}
	ok, err = session.FindOne(user)
	if err != nil || !ok {
		return protocol.Response(c, constants.ErrUserNotFound, nil)
	}

	amount := req.Amount * constants.TokenRatio
	walletType := "RMB"
	walletAddress := ""
	payChannel := "system"

	// 充值事务操作
	transactionSession := databases.NewSessionWrapper(u.dao)
	tsErr := transactionSession.Execute(func(session databases.SessionDao) error {
		payUserWallet := &models.UserWallet{UserId: user.Id}
		ok, err = session.FindOne(payUserWallet)
		if err != nil {
			return err
		}
		now := time.Now().Unix()
		if !ok {
			// 用户钱包不存在，创建一个
			payUserWallet.UserId = user.Id
			payUserWallet.WalletType = walletType
			payUserWallet.WalletAddress = walletAddress
			payUserWallet.Balance = amount
			payUserWallet.CreatedAt = now
			payUserWallet.UpdatedAt = now
			if _, err = session.InsertOne(payUserWallet); err != nil {
				u.logger.Error("insert user wallet fail",
					zap.Int64("adminUserId", adminUserId),
					zap.String("userName", user.UserName),
					zap.Int64("payUserId", user.Id),
					zap.Int64("amount", amount),
					zap.String("reason", req.Reason))
				return err
			}
		} else {
			// 钱包存在，更新余额
			payUserWallet.UpdatedAt = now
			payUserWallet.Balance += amount
			if _, err = session.UpdateById(payUserWallet.Id, payUserWallet); err != nil {
				u.logger.Error("update user wallet fail",
					zap.Int64("adminUserId", adminUserId),
					zap.Int64("payUserId", user.Id),
					zap.Int64("amount", amount),
					zap.String("reason", req.Reason))
				return err
			}
		}
		return nil
	}).Execute(func(session databases.SessionDao) error {
		// 增加充值日志
		payLog := &models.UserPayLog{
			UserId:      user.Id,
			AdminUserId: adminUserId,
			PayAmount:   amount,
			PayReason:   req.Reason,
			PayChannel:  payChannel,
			PayTime:     time.Now().Unix(),
		}
		if _, err = session.InsertOne(payLog); err != nil {
			u.logger.Error("insert user pay log fail",
				zap.Int64("adminUserId", adminUserId),
				zap.Int64("payUserId", user.Id),
				zap.Int64("amount", amount),
				zap.String("userName", user.UserName),
				zap.String("reason", req.Reason))
			return err
		}
		return nil
	}).CommitAndClose()

	if tsErr != nil {
		u.logger.Error("用户充值事务失败", zap.Error(tsErr))
		return protocol.Response(c, constants.ErrUserPayment.AppendErrors(tsErr), nil)
	}

	u.logger.Info("用户充值成功",
		zap.Int64("adminUserId", adminUserId),
		zap.Int64("userId", user.Id),
		zap.Int64("amount", amount),
		zap.String("reason", req.Reason))

	return protocol.Response(c, nil, "用户充值完成")
}

// checkUserIsAdmin 检查用户是否是管理员
func (u *UserHttpService) checkUserIsAdmin(userId int64) (bool, error) {
	session := u.dao.NewSession()
	defer session.Close()
	adminUser := &models.PlatformUsers{}
	ok, err := session.FindById(userId, adminUser)
	if err != nil {
		return false, err
	}
	if !ok {
		return false, nil
	}
	// 检查用户是否为管理员
	return true, nil
}
