package user_service

import (
	"github.com/labstack/echo/v4"
	"github.com/stardustagi/TopLib/libs/databases"
	"github.com/stardustagi/TopLib/protocol"
	"github.com/stardustagi/TopModelsPlatform/constants"
	"github.com/stardustagi/TopModelsPlatform/models"
	"github.com/stardustagi/TopModelsPlatform/protocol/requests"
	"github.com/stardustagi/TopModelsPlatform/protocol/responses"
	"go.uber.org/zap"
)

// GetUserConsumeRecord 获取用户消费记录（管理端，查询所有用户）
// @Summary 获取用户消费记录
// @Description 获取用户消费记录，支持按用户ID、时间范围、模型等条件过滤
// @Tags User
// @Accept json
// @Produce json
// @Param request body requests.GetUserConsumeRecordReq true "获取消费记录请求"
// @Success 200 {object} responses.DefaultResponse
// @Router /user/getUserConsumeRecord [post]
func (u *UserHttpService) GetUserConsumeRecord(
	ctx echo.Context,
	req requests.GetUserConsumeRecordReq,
	resp responses.DefaultResponse,
) error {
	u.logger.Info("获取用户消费记录",
		zap.Int64("userId", req.UserId),
		zap.String("model", req.Model))

	session := u.dao.NewSession()
	defer session.Close()

	// 默认分页
	if req.PageInfo.Limit <= 0 {
		req.PageInfo.Limit = 20
	}
	if req.PageInfo.Sort == "" {
		req.PageInfo.Sort = "created_at desc"
	}

	// 构建查询条件
	query := session.Native().NewSession()
	defer query.Close()

	// 可选条件：用户ID
	if req.UserId > 0 {
		query = query.Where("user_id = ?", req.UserId)
	}

	// 可选条件：时间范围
	if req.StartTime > 0 {
		query = query.And("created_at >= ?", req.StartTime)
	}
	if req.EndTime > 0 {
		query = query.And("created_at <= ?", req.EndTime)
	}

	// 可选条件：模型
	if req.Model != "" {
		query = query.And("model = ?", req.Model)
	}

	// 可选条件：供应商ID
	if req.ProviderId != "" {
		query = query.And("actual_provider_id = ?", req.ProviderId)
	}

	// 使用 FindAndCount 查询数据和总数
	var records []models.UserConsumeRecord
	total, err := query.
		OrderBy(req.PageInfo.Sort).
		Limit(req.PageInfo.Limit, req.PageInfo.Skip).
		FindAndCount(&records)
	if err != nil {
		u.logger.Error("查询用户消费记录失败", zap.Error(err))
		return protocol.Response(ctx, constants.ErrInternalServer.AppendErrors(err), nil)
	}

	return protocol.Response(ctx, nil, map[string]interface{}{
		"records": records,
		"total":   total,
	})
}

// GetUserConsumeDetail 获取用户消费详细记录（管理端，不验证用户ID）
// @Summary 获取用户消费详细记录
// @Description 根据消费记录ID获取详细信息
// @Tags User
// @Accept json
// @Produce json
// @Param request body requests.UserConsumeDetailReq true "获取消费详情请求"
// @Success 200 {object} responses.DefaultResponse
// @Router /user/getUserConsumeDetail [post]
func (u *UserHttpService) GetUserConsumeDetail(
	ctx echo.Context,
	req requests.UserConsumeDetailReq,
	resp responses.DefaultResponse,
) error {
	u.logger.Info("获取用户消费详细记录",
		zap.Int64("consumeId", req.ConsumeId),
		zap.String("consumeType", req.ConsumeType))

	session := u.dao.NewSession()
	defer session.Close()

	// 根据 ConsumeId 查询消费记录
	consumeRecord := &models.UserConsumeRecord{Id: req.ConsumeId}
	ok, err := session.FindOne(consumeRecord)
	if err != nil {
		u.logger.Error("查询消费记录失败", zap.Error(err))
		return protocol.Response(ctx, constants.ErrInternalServer.AppendErrors(err), nil)
	}
	if !ok {
		return protocol.Response(ctx, constants.ErrNotDataSet, nil)
	}

	// 默认分页
	if req.Page.Limit <= 0 {
		req.Page.Limit = 20
	}
	if req.Page.Sort == "" {
		req.Page.Sort = "id desc"
	}

	// 创建分页对象
	pageable := databases.NewPageable(req.Page.Skip, req.Page.Limit, req.Page.Sort)

	// 根据 ConsumeType 查询不同的详情表
	var result interface{}
	var total int64

	consumeType := req.ConsumeType
	if consumeType == "" {
		consumeType = consumeRecord.ConsumeType
	}

	switch consumeType {
	case "text":
		var details []models.UserConsumeDetailText
		total, err = session.FindAndCount(&details, pageable, &models.UserConsumeDetailText{ConsumdId: req.ConsumeId})
		result = details
	case "image":
		var details []models.UserConsumeDetailImage
		total, err = session.FindAndCount(&details, pageable, &models.UserConsumeDetailImage{ConsumeId: req.ConsumeId})
		result = details
	case "video":
		var details []models.UserConsumeDetailVideo
		total, err = session.FindAndCount(&details, pageable, &models.UserConsumeDetailVideo{ConsumdId: req.ConsumeId})
		result = details
	default:
		return protocol.Response(ctx, constants.ErrInvalidParams, "不支持的消费类型")
	}

	if err != nil {
		u.logger.Error("查询消费详情失败", zap.Error(err), zap.String("consumeType", consumeType))
		return protocol.Response(ctx, constants.ErrInternalServer.AppendErrors(err), nil)
	}

	return protocol.Response(ctx, nil, map[string]interface{}{
		"consume_record": consumeRecord,
		"details":        result,
		"total":          total,
	})
}
