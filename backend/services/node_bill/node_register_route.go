package node_bill

import (
	"github.com/stardustagi/TopLib/libs/server"
	"github.com/stardustagi/TopModelsPlatform/backend"
	"github.com/stardustagi/TopModelsPlatform/protocol/requests"
	"github.com/stardustagi/TopModelsPlatform/protocol/responses"
)

func (n *NodeHttpBillService) initialization() {

	n.app.AddGroup("bill", server.Request(), backend.PlatformUserAccess())

	// DiscountRule 路由
	n.app.AddPostHandler("bill", server.NewHandler[requests.CreateDiscountRuleReq, responses.DefaultResponse](
		"createDiscountRule",
		[]string{"DiscountRule"},
		n.CreateDiscountRule,
	))
	n.app.AddPostHandler("bill", server.NewHandler[requests.UpdateDiscountRuleReq, responses.DefaultResponse](
		"updateDiscountRule",
		[]string{"DiscountRule"},
		n.UpdateDiscountRule,
	))
	n.app.AddPostHandler("bill", server.NewHandler[requests.GetDiscountRuleListReq, responses.GetDiscountRuleListResp](
		"getDiscountRuleList",
		[]string{"DiscountRule"},
		n.GetDiscountRuleList,
	))

	// UserDiscount 路由
	n.app.AddPostHandler("bill", server.NewHandler[requests.CreateUserDiscountReq, responses.DefaultResponse](
		"createUserDiscount",
		[]string{"UserDiscount"},
		n.CreateUserDiscount,
	))
	n.app.AddPostHandler("bill", server.NewHandler[requests.UpdateUserDiscountReq, responses.DefaultResponse](
		"updateUserDiscount",
		[]string{"UserDiscount"},
		n.UpdateUserDiscount,
	))
	n.app.AddPostHandler("bill", server.NewHandler[requests.GetUserDiscountListReq, responses.GetUserDiscountListResp](
		"getUserDiscountList",
		[]string{"UserDiscount"},
		n.GetUserDiscountList,
	))
}
