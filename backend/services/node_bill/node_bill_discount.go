package node_bill

import (
	"fmt"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stardustagi/TopLib/protocol"
	"github.com/stardustagi/TopModelsPlatform/constants"
	"github.com/stardustagi/TopModelsPlatform/models"
	"github.com/stardustagi/TopModelsPlatform/protocol/requests"
	"github.com/stardustagi/TopModelsPlatform/protocol/responses"
	"go.uber.org/zap"
)

// ===================== DiscountRule 接口 =====================

// CreateDiscountRule 创建折扣规则
// @Summary 创建折扣规则
// @Description 创建新的折扣规则
// @Tags DiscountRule
// @Accept json
// @Produce json
// @Param request body requests.CreateDiscountRuleReq true "创建折扣规则请求"
// @Success 200 {object} responses.DefaultResponse
// @Router /node/bill/createDiscountRule [post]
func (n *NodeHttpBillService) CreateDiscountRule(ctx echo.Context,
	req requests.CreateDiscountRuleReq, resp responses.DefaultResponse) error {
	n.logger.Info("创建折扣规则", zap.String("name", req.Name))

	session := n.dao.NewSession()
	defer session.Close()

	now := time.Now().Unix()

	// 处理 Status：nil时默认启用(1)
	status := 1
	if req.Status != nil {
		status = *req.Status
	}

	rule := &models.DiscountRule{
		Name:        req.Name,
		Description: req.Description,
		//DiscountRate: req.DiscountRate,
		Status: status,
		//CreatedAt:    now,
		//UpdatedAt:    now,
	}

	ok, err := session.FindOne(rule)
	if err != nil {
		n.logger.Error("折扣查询错误", zap.Error(err))
		return protocol.Response(ctx, constants.ErrInternalServer.AppendErrors(err), nil)
	}
	if ok {
		err = fmt.Errorf("折扣已经存在,%s,%s", rule.Name, rule.Description)
		n.logger.Error("折扣添加错误", zap.Error(err))
		return protocol.Response(ctx, constants.ErrInternalServer.AppendErrors(err), nil)
	}

	rule.CreatedAt = now
	rule.UpdatedAt = now
	rule.DiscountRate = req.DiscountRate
	_, err = session.InsertOne(rule)
	if err != nil {
		n.logger.Error("创建折扣规则失败", zap.Error(err))
		return protocol.Response(ctx, constants.ErrInternalServer.AppendErrors(err), nil)
	}

	n.logger.Info("创建折扣规则成功", zap.Int64("id", rule.Id))

	return protocol.Response(ctx, nil, map[string]interface{}{
		"id":      rule.Id,
		"message": "创建折扣规则成功",
	})
}

// UpdateDiscountRule 更新折扣规则
// @Summary 更新折扣规则
// @Description 更新已有的折扣规则
// @Tags DiscountRule
// @Accept json
// @Produce json
// @Param request body requests.UpdateDiscountRuleReq true "更新折扣规则请求"
// @Success 200 {object} responses.DefaultResponse
// @Router /node/bill/updateDiscountRule [post]
func (n *NodeHttpBillService) UpdateDiscountRule(ctx echo.Context,
	req requests.UpdateDiscountRuleReq, resp responses.DefaultResponse) error {
	n.logger.Info("更新折扣规则", zap.Int64("id", req.Id))

	session := n.dao.NewSession()
	defer session.Close()

	// 查询原有配置是否存在
	existingRule := &models.DiscountRule{Id: req.Id}
	ok, err := session.FindOne(existingRule)
	if err != nil {
		n.logger.Error("查询折扣规则失败", zap.Error(err))
		return protocol.Response(ctx, constants.ErrInternalServer.AppendErrors(err), nil)
	}
	if !ok {
		return protocol.Response(ctx, constants.ErrNotDataSet, nil)
	}

	// 更新字段
	if req.Name != "" {
		existingRule.Name = req.Name
	}
	if req.Description != "" {
		existingRule.Description = req.Description
	}
	if req.DiscountRate != nil {
		existingRule.DiscountRate = *req.DiscountRate
	}
	if req.Status != nil {
		existingRule.Status = *req.Status
	}
	existingRule.UpdatedAt = time.Now().Unix()

	// 使用 Cols 强制更新包含零值的字段
	_, err = session.Native().
		ID(req.Id).
		Cols("name", "description", "discount_rate", "status", "updated_at").
		Update(existingRule)
	if err != nil {
		n.logger.Error("更新折扣规则失败", zap.Error(err))
		return protocol.Response(ctx, constants.ErrInternalServer.AppendErrors(err), nil)
	}
	// 新增：同步更新使用该规则的用户折扣的折扣率
	now := time.Now().Unix()
	userDiscountUpdate := &models.UserDiscount{
		DiscountRate: existingRule.DiscountRate,
		UpdatedAt:    now,
	}
	_, err = session.Native().
		Where("rule_id = ?", req.Id).
		Cols("discount_rate", "updated_at").
		Update(userDiscountUpdate)
	if err != nil {
		n.logger.Error("同步更新用户折扣失败", zap.Error(err))
		return protocol.Response(ctx, constants.ErrInternalServer.AppendErrors(err), nil)
	}

	n.logger.Info("更新折扣规则成功", zap.Int64("id", req.Id))

	return protocol.Response(ctx, nil, map[string]interface{}{
		"id":      req.Id,
		"message": "更新折扣规则成功",
	})
}

// GetDiscountRuleList 获取折扣规则列表
// @Summary 获取折扣规则列表
// @Description 分页获取折扣规则列表
// @Tags DiscountRule
// @Accept json
// @Produce json
// @Param request body requests.GetDiscountRuleListReq true "获取折扣规则列表请求"
// @Success 200 {object} responses.GetDiscountRuleListResp
// @Router /node/bill/getDiscountRuleList [post]
func (n *NodeHttpBillService) GetDiscountRuleList(ctx echo.Context,
	req requests.GetDiscountRuleListReq, resp responses.GetDiscountRuleListResp) error {
	n.logger.Info("获取折扣规则列表", zap.String("name", req.Name))

	session := n.dao.NewSession()
	defer session.Close()

	// 默认分页
	if req.PageInfo.Limit <= 0 {
		req.PageInfo.Limit = constants.DefaultPageSize
	}
	if req.PageInfo.Sort == "" {
		req.PageInfo.Sort = "id desc"
	}

	// 构建查询条件
	query := session.Native().NewSession()
	defer query.Close()

	// 可选条件：名称（模糊查询）
	if req.Name != "" {
		query = query.Where("name LIKE ?", "%"+req.Name+"%")
	}

	// 可选条件：折扣率
	if req.DiscountRate != nil {
		query = query.And("discount_rate = ?", *req.DiscountRate)
	}

	// 可选条件：状态
	if req.Status != nil {
		query = query.And("status = ?", *req.Status)
	}

	// 可选条件：创建时间
	if req.CreatedAt > 0 {
		query = query.And("created_at >= ?", req.CreatedAt)
	}

	// 查询数据列表
	var rules []models.DiscountRule
	total, err := query.
		OrderBy(req.PageInfo.Sort).
		Limit(req.PageInfo.Limit, req.PageInfo.Skip).
		FindAndCount(&rules)
	if err != nil {
		n.logger.Error("查询折扣规则列表失败", zap.Error(err))
		return protocol.Response(ctx, constants.ErrInternalServer.AppendErrors(err), nil)
	}

	resp.List = rules
	resp.Total = int(total)

	return protocol.Response(ctx, nil, resp)
}

// ===================== UserDiscount 接口 =====================

// CreateUserDiscount 创建用户折扣
// @Summary 创建用户折扣
// @Description 创建新的用户折扣
// @Tags UserDiscount
// @Accept json
// @Produce json
// @Param request body requests.CreateUserDiscountReq true "创建用户折扣请求"
// @Success 200 {object} responses.DefaultResponse
// @Router /node/bill/createUserDiscount [post]
func (n *NodeHttpBillService) CreateUserDiscount(ctx echo.Context,
	req requests.CreateUserDiscountReq, resp responses.DefaultResponse) error {
	n.logger.Info("创建用户折扣",
		zap.Int64("userId", req.UserId),
	)

	session := n.dao.NewSession()
	defer session.Close()
	ruleInfo := &models.DiscountRule{Id: req.RuleId}
	ok, err := session.FindOne(ruleInfo)
	if err != nil {
		n.logger.Error("查询折扣规则失败", zap.Error(err))
		return protocol.Response(ctx, constants.ErrInternalServer.AppendErrors(err), nil)
	}
	if !ok {
		n.logger.Error("折扣规则不存在", zap.Int64("ruleId", req.RuleId))
		return protocol.Response(ctx, constants.ErrNotDataSet, nil)
	}
	now := time.Now().Unix()

	discount := &models.UserDiscount{
		UserId:       req.UserId,
		ModelId:      req.ModelId,
		RuleId:       req.RuleId,
		DiscountRate: ruleInfo.DiscountRate,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	_, err = session.InsertOne(discount)
	if err != nil {
		n.logger.Error("创建用户折扣失败", zap.Error(err))
		return protocol.Response(ctx, constants.ErrInternalServer.AppendErrors(err), nil)
	}

	n.logger.Info("创建用户折扣成功", zap.Int64("id", discount.Id))

	return protocol.Response(ctx, nil, map[string]interface{}{
		"id":      discount.Id,
		"message": "创建用户折扣成功",
	})
}

// UpdateUserDiscount 更新用户折扣
// @Summary 更新用户折扣
// @Description 更新已有的用户折扣
// @Tags UserDiscount
// @Accept json
// @Produce json
// @Param request body requests.UpdateUserDiscountReq true "更新用户折扣请求"
// @Success 200 {object} responses.DefaultResponse
// @Router /node/bill/updateUserDiscount [post]
func (n *NodeHttpBillService) UpdateUserDiscount(ctx echo.Context,
	req requests.UpdateUserDiscountReq, resp responses.DefaultResponse) error {
	n.logger.Info("更新用户折扣", zap.Int64("id", req.Id))

	session := n.dao.NewSession()
	defer session.Close()

	// 查询原有配置是否存在
	existingDiscount := &models.UserDiscount{Id: req.Id}
	ok, err := session.FindOne(existingDiscount)
	if err != nil {
		n.logger.Error("查询用户折扣失败", zap.Error(err))
		return protocol.Response(ctx, constants.ErrInternalServer.AppendErrors(err), nil)
	}
	if !ok {
		return protocol.Response(ctx, constants.ErrNotDataSet, nil)
	}

	// 更新字段
	if req.RuleId > 0 {
		existingDiscount.RuleId = req.RuleId
	} else {
		n.logger.Error("更新用户折扣失败，rule_id 必须大于0")
		return protocol.Response(ctx, constants.ErrInvalidParams.AppendErrors(fmt.Errorf("rule_id 必须大于0")), nil)
	}
	existingDiscount.UpdatedAt = time.Now().Unix()

	// 使用 Cols 强制更新包含零值的字段
	_, err = session.Native().
		ID(req.Id).
		Cols("user_id", "model_id", "rule_id", "discount_rate", "updated_at").
		Update(existingDiscount)
	if err != nil {
		n.logger.Error("更新用户折扣失败", zap.Error(err))
		return protocol.Response(ctx, constants.ErrInternalServer.AppendErrors(err), nil)
	}

	n.logger.Info("更新用户折扣成功", zap.Int64("id", req.Id))

	return protocol.Response(ctx, nil, map[string]interface{}{
		"id":      req.Id,
		"message": "更新用户折扣成功",
	})
}

// GetUserDiscountList 获取用户折扣列表
// @Summary 获取用户折扣列表
// @Description 分页获取用户折扣列表
// @Tags UserDiscount
// @Accept json
// @Produce json
// @Param request body requests.GetUserDiscountListReq true "获取用户折扣列表请求"
// @Success 200 {object} responses.GetUserDiscountListResp
// @Router /node/bill/getUserDiscountList [post]
func (n *NodeHttpBillService) GetUserDiscountList(ctx echo.Context,
	req requests.GetUserDiscountListReq, resp responses.GetUserDiscountListResp) error {
	n.logger.Info("获取用户折扣列表",
		zap.Int64("userId", req.UserId),
		zap.Int("modelId", req.ModelId))

	session := n.dao.NewSession()
	defer session.Close()

	// 默认分页
	if req.PageInfo.Limit <= 0 {
		req.PageInfo.Limit = 20
	}
	if req.PageInfo.Sort == "" {
		req.PageInfo.Sort = "id desc"
	}

	// 构建查询条件
	query := session.Native().NewSession()
	defer query.Close()

	// 可选条件：用户ID
	if req.UserId > 0 {
		query = query.Where("user_id = ?", req.UserId)
	}

	// 可选条件：模型ID
	if req.ModelId > 0 {
		query = query.And("model_id = ?", req.ModelId)
	}

	// 可选条件：规则ID
	if req.RuleId > 0 {
		query = query.And("rule_id = ?", req.RuleId)
	}

	// 可选条件：折扣率
	if req.DiscountRate != nil {
		query = query.And("discount_rate = ?", *req.DiscountRate)
	}

	// 可选条件：创建时间
	if req.CreatedAt > 0 {
		query = query.And("created_at >= ?", req.CreatedAt)
	}

	// 查询数据列表
	var discounts []models.UserDiscount
	total, err := query.
		OrderBy(req.PageInfo.Sort).
		Limit(req.PageInfo.Limit, req.PageInfo.Skip).
		FindAndCount(&discounts)
	if err != nil {
		n.logger.Error("查询用户折扣列表失败", zap.Error(err))
		return protocol.Response(ctx, constants.ErrInternalServer.AppendErrors(err), nil)
	}

	resp.List = discounts
	resp.Total = int(total)

	return protocol.Response(ctx, nil, resp)
}
