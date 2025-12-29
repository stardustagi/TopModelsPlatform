package user_service

import (
	"github.com/stardustagi/TopLib/libs/server"
	"github.com/stardustagi/TopModelsPlatform/backend"
	"github.com/stardustagi/TopModelsPlatform/protocol/requests"
	"github.com/stardustagi/TopModelsPlatform/protocol/responses"
)

func (u *UserHttpService) initialization() {
	u.app.AddGroup("user", backend.PlatformUserAccess(), server.Request())

	u.app.AddGetHandler("user", server.NewHandler(
		"list",
		[]string{"User"},
		u.List,
	))

	u.app.AddGetHandler("user", server.NewHandler[requests.UserIdReq, responses.DefaultResponse](
		"get",
		[]string{"User"},
		u.GetById,
	))

	u.app.AddPostHandler("user", server.NewHandler[requests.CreateUserReq, responses.DefaultResponse](
		"create",
		[]string{"User"},
		u.Create,
	))

	u.app.AddPostHandler("user", server.NewHandler[requests.UpdateUserReq, responses.DefaultResponse](
		"update",
		[]string{"User"},
		u.Update,
	))

	u.app.AddPostHandler("user", server.NewHandler[requests.UserIdReq, responses.DefaultResponse](
		"delete",
		[]string{"User"},
		u.Delete,
	))

	u.app.AddPostHandler("user", server.NewHandler[requests.SetActiveReq, responses.DefaultResponse](
		"setActive",
		[]string{"User"},
		u.SetActive,
	))

	u.app.AddPostHandler("user", server.NewHandler[requests.AdminAddUserReq, responses.DefaultResponse](
		"adminAddUserByName",
		[]string{"admin", "addUserByName"},
		u.AdminAddUserByName,
	))

	u.app.AddPostHandler("user", server.NewHandler[requests.UserAdminPaymentReq, responses.DefaultResponse](
		"userAdminPayment",
		[]string{"user", "payment"},
		u.UserAdminPayment,
	))
	u.app.AddPostHandler("user", server.NewHandler[requests.CreateUserRebateConfigReq, responses.DefaultResponse](
		"createUserRebateConfig",
		[]string{"user", "rebate"},
		u.CreateUserRebateConfig,
	))
	u.app.AddPostHandler("user", server.NewHandler[requests.UpdateUserRebateConfigReq, responses.DefaultResponse](
		"updateUserRebateConfig",
		[]string{"user", "rebate"},
		u.UpdateUserRebateConfig,
	))
	u.app.AddPostHandler("user", server.NewHandler[requests.GetUserRebateInfoReq, responses.GetUserRebateInfoResp](
		"getUserRebateInfo",
		[]string{"user", "rebate"},
		u.GetUserRebateInfo,
	))
	u.app.AddPostHandler("user", server.NewHandler[requests.GetModelsProviderInfoListReq, responses.GetModelsProviderInfoListResp](
		"getModelsProviderInfoListResp",
		[]string{"user", "modelProviderInfo"},
		u.GetModelsProviderInfoList,
	))
}
