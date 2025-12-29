package platform

import (
	"github.com/stardustagi/TopLib/libs/server"
	"github.com/stardustagi/TopModelsPlatform/backend"
	"github.com/stardustagi/TopModelsPlatform/protocol/requests"
	"github.com/stardustagi/TopModelsPlatform/protocol/responses"
)

func (p *PlatfromHttpService) initialization() {
	// 公开路由 - 无需认证
	p.app.AddGroup("platform", server.Request())
	p.app.AddPostHandler("platform", server.NewHandler[requests.LoginReq, responses.DefaultResponse](
		"login",
		[]string{"Platform"},
		p.Login,
	))
	p.app.AddPostHandler("platform", server.NewHandler[requests.RegisterReq, responses.DefaultResponse](
		"register",
		[]string{"Platform"},
		p.RegisterUser,
	))

	// 需要认证的路由
	p.app.AddGroup("node", backend.PlatformUserAccess(), server.Request())
	p.app.AddPostHandler("node", server.NewHandler(
		"register",
		[]string{"Node"},
		p.Register,
	))

	// 模型管理路由 - 需要认证
	p.app.AddGroup("model", backend.PlatformUserAccess(), server.Request())
	p.app.AddPostHandler("model", server.NewHandler[requests.PageReq, responses.DefaultResponse](
		"list",
		[]string{"Model"},
		p.ListModels,
	))
	p.app.AddPostHandler("model", server.NewHandler[requests.ModelIdReq, responses.DefaultResponse](
		"get",
		[]string{"Model"},
		p.GetModel,
	))
	p.app.AddPostHandler("model", server.NewHandler[requests.CreateModelReq, responses.DefaultResponse](
		"create",
		[]string{"Model"},
		p.CreateModel,
	))
	p.app.AddPostHandler("model", server.NewHandler[requests.UpdateModelReq, responses.DefaultResponse](
		"update",
		[]string{"Model"},
		p.UpdateModel,
	))
	p.app.AddPostHandler("model", server.NewHandler[requests.ModelIdReq, responses.DefaultResponse](
		"delete",
		[]string{"Model"},
		p.DeleteModel,
	))
	p.app.AddPostHandler("platform", server.NewHandler(
		"graphCode",
		[]string{"graphCode"},
		p.GetGraphCode))

	p.app.AddPostHandler("platform", server.NewHandler(
		"phoneCode",
		[]string{"phoneCode"},
		p.GetPhoneCode))
}
