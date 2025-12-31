package node_report

import (
	"encoding/json"

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
		req.Page.Sort = "id desc"
	}

	// 调用存储过程 GetUserDayReportList
	results, err := session.CallProcedure("GetUserDayReportList",
		req.UserId,
		req.Page.Skip,
		req.Page.Limit,
		req.Page.Sort,
	)
	if err != nil {
		rep.logger.Error("调用存储过程失败", zap.Error(err))
		return protocol.Response(ctx, constants.ErrInternalServer.AppendErrors(err), nil)
	}

	// 解析结果
	// 第一个结果集：数据列表
	// 第二个结果集：总数

	if len(results) >= 1 {
		// 解析数据列表
		for _, row := range results[0] {
			item := responses.DayReportItem{}
			jsonBytes, err := json.Marshal(row)
			if err != nil {
				rep.logger.Error("Marshal 失败", zap.Error(err))
				continue
			}
			err = json.Unmarshal(jsonBytes, &item)
			if err != nil {
				rep.logger.Error("Unmarshal 失败", zap.Error(err))
				continue
			}
			resp.Items = append(resp.Items, item)
		}
	}

	if len(results) >= 2 && len(results[1]) >= 1 {
		if total, ok := results[1][0]["total"].(int); ok {
			resp.Total = total
		} else {
			resp.Total = 0
		}
	}
	return protocol.Response(ctx, nil, resp)
}
