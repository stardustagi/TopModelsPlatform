package node_report

import (
	"encoding/json"
	"fmt"
	"time"

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
		req.Page.Limit = constants.DefaultPageSize
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

// GetDayReportSummary 获取调用统计汇总
// @Summary 获取调用统计汇总
// @Description 获取用户调用统计汇总报表
// @Tags Report
// @Accept json
// @Produce json
// @Param request body requests.GetDayReportSummaryReq true "获取调用统计汇总请求"
// @Success 200 {object} responses.DefaultResponse
// @Router /report/getDayReportSummary [post]
func (rep *NodeHttpReportService) GetDayReportSummary(ctx echo.Context,
	req requests.GetDayReportSummaryReq,
	resp responses.GetDayReportSummaryResp,
) error {
	rep.logger.Info("获取调用统计汇总", zap.Int64("userId", req.UserId))

	session := rep.dao.NewSession()
	defer session.Close()

	// 调用存储过程 GetUserCallStatusReport
	// 如果 UserId 为 0，传 NULL 查询所有用户
	var userId interface{}
	if req.UserId > 0 {
		userId = req.UserId
	} else {
		userId = nil
	}

	results, err := session.CallProcedure("GetUserCallStatusReport", userId)
	if err != nil {
		rep.logger.Error("调用存储过程失败", zap.Error(err))
		return protocol.Response(ctx, constants.ErrInternalServer.AppendErrors(err), nil)
	}

	// 解析结果：使用 JSON 序列化方式转换
	var list []responses.UserCallStatusItem
	if len(results) > 0 && len(results[0]) > 0 {
		// 将 results[0] 转换为 JSON bytes
		jsonBytes, err := json.Marshal(results[0])
		if err != nil {
			rep.logger.Error("序列化结果失败", zap.Error(err))
			return protocol.Response(ctx, constants.ErrInternalServer.AppendErrors(err), nil)
		}

		// 反序列化为 []responses.UserCallStatusItem
		if err := json.Unmarshal(jsonBytes, &list); err != nil {
			rep.logger.Error("反序列化结果失败", zap.Error(err))
			return protocol.Response(ctx, constants.ErrInternalServer.AppendErrors(err), nil)
		}
	}

	resp.List = list
	resp.Total = len(list)

	return protocol.Response(ctx, nil, resp)
}

func (rep *NodeHttpReportService) GetReportDetailList(ctx echo.Context,
	req requests.GetReportDetailListReq,
	resp responses.GetReportDetailListResp) error {
	rep.logger.Info("获取报表详细列表", zap.Any("req", req))

	session := rep.dao.NewSession()
	defer session.Close()

	// 默认分页
	if req.PageInfo.Limit <= 0 {
		req.PageInfo.Limit = constants.DefaultPageSize
	}
	if req.PageInfo.Sort == "" {
		req.PageInfo.Sort = "sr.id desc"
	}

	// 根据日期范围生成分表名列表
	tableNames := rep.getShardedTableNames(req.StartDate, req.EndDate)
	if len(tableNames) == 0 {
		// 没有有效的表，返回空结果
		resp.Items = []responses.GetReportDetailListItem{}
		resp.Total = 0
		return protocol.Response(ctx, nil, resp)
	}

	// 构建 UNION ALL 查询
	var allItems []responses.GetReportDetailListItem
	var totalCount int64

	for _, tableName := range tableNames {
		// 检查表是否存在
		exists, err := session.Native().IsTableExist(tableName)
		if err != nil {
			rep.logger.Warn("检查表是否存在失败", zap.Error(err), zap.String("table", tableName))
			continue
		}
		if !exists {
			rep.logger.Debug("表不存在，跳过", zap.String("table", tableName))
			continue
		}

		// 构建查询条件，关联用户表获取用户名
		query := session.Native().
			Table(tableName).Alias("sr").
			Join("LEFT", []string{"users", "u"}, "sr.user_id = u.id").
			Select("sr.id, sr.model, sr.model_id, sr.actual_model, sr.provider, sr.actual_provider, sr.actual_provider_id, sr.caller_key, sr.stream, sr.report_type, sr.tokens_per_sec, sr.latency, sr.step, sr.status_code, sr.status_message, sr.created_at, sr.user_id, u.user_name")

		// 必选条件：时间范围
		if req.StartDate > 0 {
			startTime := time.Unix(req.StartDate, 0)
			query = query.Where("sr.created_at >= ?", startTime)
		}
		if req.EndDate > 0 {
			endTime := time.Unix(req.EndDate, 0)
			query = query.And("sr.created_at <= ?", endTime)
		}

		// 可选条件：CallerKey
		if req.CallKey != "" {
			query = query.And("sr.caller_key = ?", req.CallKey)
		}

		// 可选条件：Model（模糊查询）
		if req.Model != "" {
			query = query.And("sr.model LIKE ?", "%"+req.Model+"%")
		}

		// 查询当前表的数据
		var items []responses.GetReportDetailListItem
		count, err := query.
			OrderBy(req.PageInfo.Sort).
			FindAndCount(&items)
		if err != nil {
			rep.logger.Error("查询报表详细列表失败", zap.Error(err), zap.String("table", tableName))
			continue
		}

		allItems = append(allItems, items...)
		totalCount += count
	}

	// 对合并后的结果进行分页
	start := req.PageInfo.Skip
	end := req.PageInfo.Skip + req.PageInfo.Limit
	if start > len(allItems) {
		start = len(allItems)
	}
	if end > len(allItems) {
		end = len(allItems)
	}

	resp.Items = allItems[start:end]
	resp.Total = int(totalCount)

	return protocol.Response(ctx, nil, resp)
}

// getShardedTableNames 根据日期范围生成分表名列表
func (rep *NodeHttpReportService) getShardedTableNames(startDate, endDate int64) []string {
	var tableNames []string

	if startDate <= 0 || endDate <= 0 {
		// 如果没有指定日期，使用当天的表
		t := time.Now()
		tableName := fmt.Sprintf("status_report_%d%02d%02d", t.Year(), t.Month(), t.Day())
		return []string{tableName}
	}

	startTime := time.Unix(startDate, 0)
	endTime := time.Unix(endDate, 0)

	// 遍历日期范围，生成每天的表名
	for current := startTime; !current.After(endTime); current = current.AddDate(0, 0, 1) {
		tableName := fmt.Sprintf("status_report_%d%02d%02d", current.Year(), current.Month(), current.Day())
		tableNames = append(tableNames, tableName)
	}
	return tableNames
}
