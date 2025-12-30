package platform

import (
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

// ListModels 获取模型列表
// @Summary 获取模型列表
// @Description 分页获取模型列表
// @Tags Model
// @Accept json
// @Produce json
// @Param skip query int false "跳过条数"
// @Param limit query int false "每页条数"
// @Success 200 {object} responses.DefaultResponse
// @Router /model/list [get]
func (p *PlatfromHttpService) ListModels(c echo.Context, req requests.PageReq, resp responses.DefaultResponse) error {
	p.logger.Info("获取模型列表", zap.Int("skip", req.Skip), zap.Int("limit", req.Limit))

	session := p.dao.NewSession()
	defer session.Close()

	if req.Limit <= 0 {
		req.Limit = 20
	}
	if req.Limit > 100 {
		req.Limit = 100
	}
	if req.Sort == "" {
		req.Sort = "id DESC"
	}

	var modelsList []models.ModelsInfo
	pageable := databases.NewPageable(req.Skip, req.Limit, req.Sort)
	total, err := session.FindAndCount(&modelsList, pageable)
	if err != nil {
		p.logger.Error("查询模型列表失败", zap.Error(err))
		return protocol.Response(c, constants.ErrInternalServer, nil)
	}

	result := map[string]interface{}{
		"list":  modelsList,
		"total": total,
		"skip":  req.Skip,
		"limit": req.Limit,
	}
	return protocol.Response(c, nil, result)
}

// GetModel 获取模型详情
// @Summary 获取模型详情
// @Description 根据ID获取模型详情
// @Tags Model
// @Accept json
// @Produce json
// @Param id query int true "模型ID"
// @Success 200 {object} responses.DefaultResponse
// @Router /model/get [get]
func (p *PlatfromHttpService) GetModel(c echo.Context, req requests.ModelIdReq, resp responses.DefaultResponse) error {
	p.logger.Info("获取模型详情", zap.Int64("modelId", req.Id))

	session := p.dao.NewSession()
	defer session.Close()

	model := &models.ModelsInfo{Id: req.Id}
	ok, err := session.FindOne(model)
	if err != nil {
		p.logger.Error("查询模型失败", zap.Error(err))
		return protocol.Response(c, constants.ErrInternalServer, nil)
	}
	if !ok {
		return protocol.Response(c, constants.ErrModelNotFound, nil)
	}

	return protocol.Response(c, nil, model)
}

// CreateModel 创建模型
// @Summary 创建模型
// @Description 创建新模型
// @Tags Model
// @Accept json
// @Produce json
// @Param request body requests.CreateModelReq true "创建模型请求"
// @Success 200 {object} responses.DefaultResponse
// @Router /model/create [post]
func (p *PlatfromHttpService) CreateModel(c echo.Context, req requests.CreateModelReq, resp responses.DefaultResponse) error {
	p.logger.Info("创建模型", zap.String("modelId", req.Name))

	session := p.dao.NewSession()
	defer session.Close()

	existModel := &models.ModelsInfo{ModelId: req.Name}
	ok, _ := session.FindOne(existModel)
	if ok {
		return protocol.Response(c, constants.ErrModelAlreadyExists, nil)
	}

	now := time.Now().Unix()
	model := &models.ModelsInfo{
		Name:        req.Name,
		ApiVersion:  req.ApiVersion,
		DeployName:  req.DeployName,
		InputPrice:  req.InputPrice,
		OutputPrice: req.OutputPrice,
		CachePrice:  req.CachePrice,
		Status:      "active",
		LastUpdate:  now,
		IsPrivate:   req.IsPrivate,
		OwnerId:     req.OwnerId,
		Address:     req.Address,
		ApiStyles:   req.ApiStyles,
	}

	_, err := session.InsertOne(model)
	if err != nil {
		p.logger.Error("创建模型失败", zap.Error(err))
		return protocol.Response(c, constants.ErrInternalServer, nil)
	}

	return protocol.Response(c, nil, model)
}

// UpdateModel 更新模型
// @Summary 更新模型
// @Description 更新模型信息
// @Tags Model
// @Accept json
// @Produce json
// @Param request body requests.UpdateModelReq true "更新模型请求"
// @Success 200 {object} responses.DefaultResponse
// @Router /model/update [post]
func (p *PlatfromHttpService) UpdateModel(c echo.Context, req requests.UpdateModelReq, resp responses.DefaultResponse) error {
	p.logger.Info("更新模型", zap.Int64("modelId", req.Id))

	session := p.dao.NewSession()
	defer session.Close()

	model := &models.ModelsInfo{Id: req.Id}
	ok, err := session.FindOne(model)
	if err != nil || !ok {
		return protocol.Response(c, constants.ErrModelNotFound, nil)
	}

	if req.Name != "" {
		model.Name = req.Name
	}
	if req.ApiVersion != "" {
		model.ApiVersion = req.ApiVersion
	}
	if req.DeployName != "" {
		model.DeployName = req.DeployName
	}
	if req.InputPrice > 0 {
		model.InputPrice = req.InputPrice
	}
	if req.OutputPrice > 0 {
		model.OutputPrice = req.OutputPrice
	}
	if req.Status != "" {
		model.Status = req.Status
	}
	if req.Address != "" {
		model.Address = req.Address
	}
	if req.ApiStyles != "" {
		model.ApiStyles = req.ApiStyles
	}
	model.LastUpdate = time.Now().Unix()

	_, err = session.UpdateById(model.Id, model)
	if err != nil {
		p.logger.Error("更新模型失败", zap.Error(err))
		return protocol.Response(c, constants.ErrInternalServer, nil)
	}

	return protocol.Response(c, nil, model)
}

// DeleteModel 删除模型
// @Summary 删除模型
// @Description 删除指定模型
// @Tags Model
// @Accept json
// @Produce json
// @Param request body requests.ModelIdReq true "模型ID"
// @Success 200 {object} responses.DefaultResponse
// @Router /model/delete [post]
func (p *PlatfromHttpService) DeleteModel(c echo.Context, req requests.ModelIdReq, resp responses.DefaultResponse) error {
	p.logger.Info("删除模型", zap.Int64("modelId", req.Id))

	session := p.dao.NewSession()
	defer session.Close()

	model := &models.ModelsInfo{Id: req.Id}
	ok, err := session.FindOne(model)
	if err != nil || !ok {
		return protocol.Response(c, constants.ErrModelNotFound, nil)
	}

	_, err = session.Delete(model)
	if err != nil {
		p.logger.Error("删除模型失败", zap.Error(err))
		return protocol.Response(c, constants.ErrInternalServer, nil)
	}

	return protocol.Response(c, nil, nil)
}
