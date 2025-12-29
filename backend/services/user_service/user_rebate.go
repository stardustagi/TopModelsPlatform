package user_service

import (
	"fmt"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stardustagi/TopLib/libs/databases"
	"github.com/stardustagi/TopLib/protocol"
	"github.com/stardustagi/TopModelsPlatform/constants"
	"github.com/stardustagi/TopModelsPlatform/models"
	"github.com/stardustagi/TopModelsPlatform/protocol/requests"
	"github.com/stardustagi/TopModelsPlatform/protocol/responses"
	"go.uber.org/zap"
)

// CreateUserRebateConfig 创建用户返利配置
// @Summary 创建用户返利配置
// @Description 创建用户消费返点阶梯配置
// @Tags User
// @Accept json
// @Produce json
// @Param request body requests.CreateUserRebateConfigReq true "创建返利配置请求"
// @Success 200 {object} responses.DefaultResponse
// @Router /user/createUserRebateConfig [post]
func (u *UserHttpService) CreateUserRebateConfig(ctx echo.Context,
	req requests.CreateUserRebateConfigReq, resp responses.DefaultResponse) error {
	u.logger.Info("创建用户返利配置",
		zap.Int64("userId", req.UserId),
		zap.Int64("tierStart", req.TierStart),
		zap.Int64("tierEnd", req.TierEnd),
		zap.Int("rebateRate", req.RebateRate))

	session := u.dao.NewSession()
	defer session.Close()

	// 验证用户是否存在
	user := &models.Users{Id: req.UserId}
	ok, err := session.FindOne(user)
	if err != nil {
		u.logger.Error("查询用户失败", zap.Error(err), zap.Int64("userId", req.UserId))
		return protocol.Response(ctx, constants.ErrInternalServer.AppendErrors(err), nil)
	}
	if !ok {
		return protocol.Response(ctx, constants.ErrUserNotFound, nil)
	}

	// 检查是否存在重叠的阶梯配置
	var existingConfigs []models.UserRebateConfig
	err = session.Native().
		Where("user_id = ? AND status = 1", req.UserId).
		Find(&existingConfigs)
	if err != nil {
		u.logger.Error("查询已有返利配置失败", zap.Error(err))
		return protocol.Response(ctx, constants.ErrInternalServer.AppendErrors(err), nil)
	}

	// 验证阶梯区间是否重叠
	for _, config := range existingConfigs {
		configEnd := config.TierEnd
		if configEnd == -1 {
			configEnd = int64(^uint64(0) >> 1)
		}
		reqEnd := req.TierEnd
		if reqEnd == -1 {
			reqEnd = int64(^uint64(0) >> 1)
		}

		if req.TierStart < configEnd && reqEnd > config.TierStart {
			u.logger.Warn("阶梯区间重叠",
				zap.Int64("existingStart", config.TierStart),
				zap.Int64("existingEnd", config.TierEnd),
				zap.Int64("newStart", req.TierStart),
				zap.Int64("newEnd", req.TierEnd))
			return protocol.Response(ctx,
				constants.ErrInvalidParams.AppendErrors(fmt.Errorf("阶梯区间与现有配置重叠")), nil)
		}
	}

	// 创建返利配置
	now := time.Now().Unix()
	rebateConfig := &models.UserRebateConfig{
		UserId:     req.UserId,
		TierStart:  req.TierStart,
		TierEnd:    req.TierEnd,
		RebateRate: req.RebateRate,
		Status:     req.Status,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	// 默认启用
	if req.Status == 0 {
		rebateConfig.Status = 1
	}

	_, err = session.InsertOne(rebateConfig)
	if err != nil {
		u.logger.Error("创建返利配置失败", zap.Error(err))
		return protocol.Response(ctx, constants.ErrInternalServer.AppendErrors(err), nil)
	}

	u.logger.Info("创建用户返利配置成功",
		zap.Int64("configId", rebateConfig.Id),
		zap.Int64("userId", req.UserId))

	return protocol.Response(ctx, nil, map[string]interface{}{
		"id":      rebateConfig.Id,
		"message": "创建用户返利配置成功",
	})
}

// UpdateUserRebateConfig 修改用户返利配置
// @Summary 修改用户返利配置
// @Description 修改用户消费返点阶梯配置
// @Tags User
// @Accept json
// @Produce json
// @Param request body requests.UpdateUserRebateConfigReq true "修改返利配置请求"
// @Success 200 {object} responses.DefaultResponse
// @Router /user/updateUserRebateConfig [post]
func (u *UserHttpService) UpdateUserRebateConfig(ctx echo.Context,
	req requests.UpdateUserRebateConfigReq, resp responses.DefaultResponse) error {
	u.logger.Info("修改用户返利配置",
		zap.Int64("id", req.Id),
		zap.Int64("userId", req.UserId),
		zap.Int64("tierStart", req.TierStart),
		zap.Int64("tierEnd", req.TierEnd),
		zap.Int("rebateRate", req.RebateRate))

	session := u.dao.NewSession()
	defer session.Close()

	// 查询原有配置是否存在
	existingConfig := &models.UserRebateConfig{Id: req.Id}
	ok, err := session.FindOne(existingConfig)
	if err != nil {
		u.logger.Error("查询返利配置失败", zap.Error(err), zap.Int64("id", req.Id))
		return protocol.Response(ctx, constants.ErrInternalServer.AppendErrors(err), nil)
	}
	if !ok {
		return protocol.Response(ctx, constants.ErrNotDataSet, nil)
	}

	// 验证用户ID是否匹配
	if existingConfig.UserId != req.UserId {
		return protocol.Response(ctx, constants.ErrInvalidParams.AppendErrors(fmt.Errorf("用户ID不匹配")), nil)
	}

	// 检查是否存在重叠的阶梯配置（排除当前配置）
	var otherConfigs []models.UserRebateConfig
	err = session.Native().
		Where("user_id = ? AND status = 1 AND id != ?", req.UserId, req.Id).
		Find(&otherConfigs)
	if err != nil {
		u.logger.Error("查询已有返利配置失败", zap.Error(err))
		return protocol.Response(ctx, constants.ErrInternalServer.AppendErrors(err), nil)
	}

	// 验证阶梯区间是否重叠
	for _, config := range otherConfigs {
		configEnd := config.TierEnd
		if configEnd == -1 {
			configEnd = int64(^uint64(0) >> 1)
		}
		reqEnd := req.TierEnd
		if reqEnd == -1 {
			reqEnd = int64(^uint64(0) >> 1)
		}

		if req.TierStart < configEnd && reqEnd > config.TierStart {
			u.logger.Warn("阶梯区间重叠",
				zap.Int64("existingStart", config.TierStart),
				zap.Int64("existingEnd", config.TierEnd),
				zap.Int64("newStart", req.TierStart),
				zap.Int64("newEnd", req.TierEnd))
			return protocol.Response(ctx,
				constants.ErrInvalidParams.AppendErrors(fmt.Errorf("阶梯区间与现有配置重叠")), nil)
		}
	}

	// 更新返利配置
	existingConfig.TierStart = req.TierStart
	existingConfig.TierEnd = req.TierEnd
	existingConfig.RebateRate = req.RebateRate
	existingConfig.Status = req.Status
	existingConfig.UpdatedAt = time.Now().Unix()

	_, err = session.UpdateById(req.Id, existingConfig)
	if err != nil {
		u.logger.Error("修改返利配置失败", zap.Error(err))
		return protocol.Response(ctx, constants.ErrInternalServer.AppendErrors(err), nil)
	}

	u.logger.Info("修改用户返利配置成功",
		zap.Int64("configId", req.Id),
		zap.Int64("userId", req.UserId))

	return protocol.Response(ctx, nil, map[string]interface{}{
		"id":      req.Id,
		"message": "修改用户返利配置成功",
	})
}

// GetUserRebateInfo 获取用户返利信息
// @Summary 获取用户返利信息
// @Description 获取用户月度返点记录，UserId为0时查询所有记录
// @Tags User
// @Accept json
// @Produce json
// @Param request body requests.GetUserRebateInfoReq true "获取返利信息请求"
// @Success 200 {object} responses.GetUserRebateInfoResp
// @Router /user/getUserRebateInfo [post]
func (u *UserHttpService) GetUserRebateInfo(ctx echo.Context,
	req requests.GetUserRebateInfoReq,
	resp responses.GetUserRebateInfoResp) error {
	u.logger.Info("获取用户返利信息",
		zap.Int64("userId", req.UserId))

	session := u.dao.NewSession()
	defer session.Close()

	// 默认分页
	if req.PageInfo.Limit <= 0 {
		req.PageInfo.Limit = 20
	}
	if req.PageInfo.Sort == "" {
		req.PageInfo.Sort = "id desc"
	}

	var records []models.UserRebateMonthly
	var total int64
	var err error

	// 构建查询条件
	queryModel := &models.UserRebateMonthly{}
	if req.UserId > 0 {
		queryModel.UserId = req.UserId
	}

	// 创建分页对象
	pageable := databases.NewPageable(req.PageInfo.Skip, req.PageInfo.Limit, req.PageInfo.Sort)

	// 使用 FindAndCount 查询
	total, err = session.FindAndCount(&records, pageable, queryModel)
	if err != nil {
		u.logger.Error("查询用户返利记录失败", zap.Error(err))
		return protocol.Response(ctx, constants.ErrInternalServer.AppendErrors(err), nil)
	}

	// 转换为响应数据
	var data []responses.UserRebateData
	for _, r := range records {
		data = append(data, responses.UserRebateData{
			UserId:        r.UserId,
			Month:         r.Month,
			TotalConsumed: r.TotalConsumed,
			RebateAmount:  r.RebateAmount,
			RebateUsed:    r.RebateUsed,
			RebateRate:    r.RebateRate,
			Status:        r.Status,
		})
	}

	resp.Data = data
	resp.Total = int(total)

	return protocol.Response(ctx, nil, resp)
}
