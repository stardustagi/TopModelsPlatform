package node_report

import (
	"github.com/labstack/echo/v4"
	"github.com/stardustagi/TopLib/protocol"
	"github.com/stardustagi/TopModelsPlatform/constants"
	"github.com/stardustagi/TopModelsPlatform/protocol/requests"
	"github.com/stardustagi/TopModelsPlatform/protocol/responses"
	"go.uber.org/zap"
)

// GetDayReportList 获取日报表列表
// @Summary 获取日报表列表
// @Description 分页获取用户日消费报表
// @Tags Report
// @Accept json
// @Produce json
// @Param request body requests.GetDayReportListReq true "获取日报表请求"
// @Success 200 {object} responses.DefaultResponse
// @Router /report/getDayReportList [post]
func (rep *NodeHttpReportService) GetDayReportList(ctx echo.Context,
	req requests.GetDayReportListReq, resp responses.GetDayReportListResp) error {
	rep.logger.Info("获取日报表列表",
		zap.Int64("userId", req.UserId),
		zap.Int("skip", req.Page.Skip),
		zap.Int("limit", req.Page.Limit))

	session := rep.dao.NewSession()
	defer session.Close()

	// 默认分页
	if req.Page.Limit <= 0 {
		req.Page.Limit = 20
	}
	if req.Page.Sort == "" {
		req.Page.Sort = "pds.id desc"
	}

	// 构建查询条件
	query := session.Native().
		Table("provider_model_daily_summary").Alias("pds").
		Join("LEFT", []string{"users", "u"}, "pds.user_id = u.id").
		Join("LEFT", []string{"models_provider", "mp"}, "pds.actual_provider_id = mp.id").
		Join("LEFT", []string{"models_info", "mi"}, "pds.model_id = mi.id").
		Select("pds.id, pds.user_id, u.user_name, pds.actual_provider_id, mp.name as provider_name, pds.model_id, mi.name as model_name, pds.consume_type, pds.date, pds.total_consumed, pds.total_cost, pds.updated_at")

	// 可选条件：用户ID
	if req.UserId > 0 {
		query = query.Where("pds.user_id = ?", req.UserId)
	}

	// 可选条件：用户名（模糊查询）
	if req.UserName != "" {
		query = query.And("u.user_name LIKE ?", "%"+req.UserName+"%")
	}

	// 可选条件：供应商名称（模糊查询）
	if req.ProviderName != "" {
		query = query.And("mp.name LIKE ?", "%"+req.ProviderName+"%")
	}

	// 可选条件：模型名称（模糊查询）
	if req.ModelName != "" {
		query = query.And("mi.name LIKE ?", "%"+req.ModelName+"%")
	}

	// 可选条件：消费类型
	if req.ConsumeType != "" {
		query = query.And("pds.consume_type = ?", req.ConsumeType)
	}

	// 可选条件：更新时间
	if req.UpdatedAt > 0 {
		query = query.And("pds.updated_at >= ?", req.UpdatedAt)
	}

	// 可选条件：总成本（大于等于）
	if req.TotalCost > 0 {
		query = query.And("pds.total_cost >= ?", req.TotalCost)
	}

	// 查询数据列表
	var items []responses.DayReportItem
	total, err := query.
		OrderBy(req.Page.Sort).
		Limit(req.Page.Limit, req.Page.Skip).
		FindAndCount(&items)
	if err != nil {
		rep.logger.Error("查询日报表列表失败", zap.Error(err))
		return protocol.Response(ctx, constants.ErrInternalServer.AppendErrors(err), nil)
	}

	resp.Items = items
	resp.Total = int(total)

	return protocol.Response(ctx, nil, resp)
}
