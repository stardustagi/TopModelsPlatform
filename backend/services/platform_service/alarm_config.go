package platform

import (
	"encoding/json"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stardustagi/TopLib/libs/redis"
	"github.com/stardustagi/TopLib/protocol"
	"github.com/stardustagi/TopModelsPlatform/constants"
	"github.com/stardustagi/TopModelsPlatform/models"
	"github.com/stardustagi/TopModelsPlatform/protocol/requests"
	"github.com/stardustagi/TopModelsPlatform/protocol/responses"
	"go.uber.org/zap"
)

// CreateAlarm 创建告警配置
// @Summary 创建告警配置
// @Description 创建新的告警配置
// @Tags Alarm
// @Accept json
// @Produce json
// @Param request body requests.CreateAlarmReq true "创建告警配置请求"
// @Success 200 {object} responses.DefaultResponse
// @Router /platform/createAlarm [post]
func (p *PlatfromHttpService) CreateAlarm(ctx echo.Context,
	req requests.CreateAlarmReq, resp responses.DefaultResponse) error {
	p.logger.Info("创建告警配置",
		zap.String("type", req.Type),
		zap.Int64("userId", *req.UserId))

	session := p.dao.NewSession()
	defer session.Close()

	now := time.Now().Unix()

	// 处理 Status：nil时默认启用(1)
	status := 1
	if req.Status != nil {
		status = *req.Status
	}
	userId := int64(0)
	if req.UserId != nil {
		userId = *req.UserId
	}
	alarmConfig := &models.AlarmConfig{
		Type:       req.Type,
		Min:        req.Min,
		Max:        req.Max,
		Status:     status,
		UserId:     *req.UserId,
		CreatedAt:  now,
		LastupDate: now,
	}

	_, err := session.InsertOne(alarmConfig)
	if err != nil {
		p.logger.Error("创建告警配置失败", zap.Error(err))
		return protocol.Response(ctx, constants.ErrInternalServer.AppendErrors(err), nil)
	}

	// 将告警配置写入Redis
	redisKey := constants.PlatformAlarmKey(userId)
	rds := redis.GetRedisDb()
	alarmJson, err := json.Marshal(alarmConfig)
	if err != nil {
		p.logger.Error("序列化告警配置失败", zap.Error(err))
		return protocol.Response(ctx, constants.ErrInternalServer.AppendErrors(err), nil)
	}
	err = rds.Set(p.ctx, redisKey, alarmJson, 0).Err()
	if err != nil {
		p.logger.Error("写入Redis告警配置失败", zap.Error(err), zap.String("key", redisKey))
		return protocol.Response(ctx, constants.ErrInternalServer.AppendErrors(err), nil)
	}

	p.logger.Info("创建告警配置成功", zap.Int64("id", alarmConfig.Id), zap.String("redisKey", redisKey))

	return protocol.Response(ctx, nil, map[string]interface{}{
		"id":      alarmConfig.Id,
		"message": "创建告警配置成功",
	})
}

// UpdateAlarm 更新告警配置
// @Summary 更新告警配置
// @Description 更新已有的告警配置
// @Tags Alarm
// @Accept json
// @Produce json
// @Param request body requests.UpdateAlarmReq true "更新告警配置请求"
// @Success 200 {object} responses.DefaultResponse
// @Router /platform/updateAlarm [post]
func (p *PlatfromHttpService) UpdateAlarm(ctx echo.Context,
	req requests.UpdateAlarmReq, resp responses.DefaultResponse) error {
	p.logger.Info("更新告警配置", zap.Int64("id", req.Id))

	session := p.dao.NewSession()
	defer session.Close()

	// 查询原有配置是否存在
	existingConfig := &models.AlarmConfig{Id: req.Id}
	ok, err := session.FindOne(existingConfig)
	if err != nil {
		p.logger.Error("查询告警配置失败", zap.Error(err))
		return protocol.Response(ctx, constants.ErrInternalServer.AppendErrors(err), nil)
	}
	if !ok {
		return protocol.Response(ctx, constants.ErrNotDataSet, nil)
	}

	// 更新字段
	if req.Type != "" {
		existingConfig.Type = req.Type
	}
	if req.Min != nil {
		existingConfig.Min = *req.Min
	}
	if req.Max != nil {
		existingConfig.Max = *req.Max
	}
	if req.Status != nil {
		existingConfig.Status = *req.Status
	}
	if req.UserId != nil {
		existingConfig.UserId = *req.UserId
	} else {
		existingConfig.UserId = 0
	}
	existingConfig.LastupDate = time.Now().Unix()

	// 使用 Cols 强制更新包含零值的字段
	_, err = session.Native().
		ID(req.Id).
		Cols("type", "min", "max", "status", "user_id", "lastup_date").
		Update(existingConfig)
	if err != nil {
		p.logger.Error("更新告警配置失败", zap.Error(err))
		return protocol.Response(ctx, constants.ErrInternalServer.AppendErrors(err), nil)
	}

	// 将更新后的告警配置写入Redis
	redisKey := constants.PlatformAlarmKey(req.Id)
	alarmJson, err := json.Marshal(existingConfig)
	if err != nil {
		p.logger.Error("序列化告警配置失败", zap.Error(err))
		return protocol.Response(ctx, constants.ErrInternalServer.AppendErrors(err), nil)
	}
	rds := redis.GetRedisDb()
	err = rds.Set(p.ctx, redisKey, alarmJson, 0).Err()
	if err != nil {
		p.logger.Error("写入Redis告警配置失败", zap.Error(err), zap.String("key", redisKey))
		return protocol.Response(ctx, constants.ErrInternalServer.AppendErrors(err), nil)
	}

	p.logger.Info("更新告警配置成功", zap.Int64("id", req.Id), zap.String("redisKey", redisKey))

	return protocol.Response(ctx, nil, map[string]interface{}{
		"id":      req.Id,
		"message": "更新告警配置成功",
	})
}

// GetAlarm 获取告警配置
// @Summary 获取告警配置
// @Description 查询告警配置列表或单个配置
// @Tags Alarm
// @Accept json
// @Produce json
// @Param request body requests.GetAlarmReq true "获取告警配置请求"
// @Success 200 {object} responses.DefaultResponse
// @Router /platform/getAlarm [post]
func (p *PlatfromHttpService) GetAlarm(ctx echo.Context,
	req requests.GetAlarmReq, resp responses.DefaultResponse) error {
	p.logger.Info("获取告警配置",
		zap.Int64("id", req.Id),
		zap.Int64("userId", req.UserId))

	session := p.dao.NewSession()
	defer session.Close()

	// 如果指定了ID，查询单条记录
	if req.Id > 0 {
		alarmConfig := &models.AlarmConfig{Id: req.Id}
		ok, err := session.FindOne(alarmConfig)
		if err != nil {
			p.logger.Error("查询告警配置失败", zap.Error(err))
			return protocol.Response(ctx, constants.ErrInternalServer.AppendErrors(err), nil)
		}
		if !ok {
			return protocol.Response(ctx, constants.ErrNotDataSet, nil)
		}
		return protocol.Response(ctx, nil, alarmConfig)
	}

	// 默认分页
	if req.Page.Limit <= 0 {
		req.Page.Limit = 20
	}
	if req.Page.Sort == "" {
		req.Page.Sort = "id desc"
	}

	// 构建查询条件
	query := session.Native().NewSession()
	defer query.Close()

	// 可选条件：用户ID
	if req.UserId > 0 {
		query = query.Where("user_id = ?", req.UserId)
	}

	// 可选条件：类型
	if req.Type != "" {
		query = query.And("type = ?", req.Type)
	}

	// 可选条件：最小值
	if req.Min != nil {
		query = query.And("min >= ?", *req.Min)
	}

	// 可选条件：最大值
	if req.Max != nil {
		query = query.And("max <= ?", *req.Max)
	}

	// 可选条件：创建时间
	if req.CreatedAt > 0 {
		query = query.And("created_at >= ?", req.CreatedAt)
	}

	// 可选条件：状态（使用指针区分未传值和0）
	if req.Status != nil {
		query = query.And("status = ?", *req.Status)
	}

	// 查询列表
	var alarmConfigs []models.AlarmConfig
	total, err := query.
		OrderBy(req.Page.Sort).
		Limit(req.Page.Limit, req.Page.Skip).
		FindAndCount(&alarmConfigs)
	if err != nil {
		p.logger.Error("查询告警配置列表失败", zap.Error(err))
		return protocol.Response(ctx, constants.ErrInternalServer.AppendErrors(err), nil)
	}

	return protocol.Response(ctx, nil, map[string]interface{}{
		"list":  alarmConfigs,
		"total": total,
	})
}
