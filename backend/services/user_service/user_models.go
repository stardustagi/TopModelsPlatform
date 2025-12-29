package user_service

import (
	"github.com/labstack/echo/v4"
	"github.com/stardustagi/TopLib/protocol"
	"github.com/stardustagi/TopModelsPlatform/constants"
	"github.com/stardustagi/TopModelsPlatform/models"
	"github.com/stardustagi/TopModelsPlatform/protocol/requests"
	"github.com/stardustagi/TopModelsPlatform/protocol/responses"
	"go.uber.org/zap"
)

// GetModelsProviderInfoList 获取模型供应商列表
// @Summary 获取模型供应商列表
// @Description 获取模型供应商列表，支持按名称模糊查询
// @Tags User
// @Accept json
// @Produce json
// @Param request body requests.GetModelsProviderInfoListReq true "获取供应商列表请求"
// @Success 200 {object} responses.GetModelsProviderInfoListResp
// @Router /user/getModelsProviderInfoList [post]
func (u *UserHttpService) GetModelsProviderInfoList(ctx echo.Context,
	req requests.GetModelsProviderInfoListReq,
	resp responses.GetModelsProviderInfoListResp) error {
	u.logger.Info("获取模型供应商列表", zap.String("providerName", req.ProviderName))

	session := u.dao.NewSession()
	defer session.Close()

	var providers []models.ModelsProvider
	var err error

	query := session.Native().Where("deleted = 0")

	// 如果name不为空，则like查询
	if req.ProviderName != "" {
		query = query.And("name LIKE ?", "%"+req.ProviderName+"%")
	}

	err = query.Find(&providers)
	if err != nil {
		u.logger.Error("查询模型供应商列表失败", zap.Error(err))
		return protocol.Response(ctx, constants.ErrInternalServer.AppendErrors(err), nil)
	}

	// 转换为响应数据
	var data []responses.ModelsProviderInfo
	for _, p := range providers {
		data = append(data, responses.ModelsProviderInfo{
			Id:   p.Id,
			Name: p.Name,
		})
	}

	resp.Providers = data
	resp.Total = len(data)

	return protocol.Response(ctx, nil, resp)
}
