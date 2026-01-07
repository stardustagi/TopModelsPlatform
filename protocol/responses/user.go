package responses

type ModelsProviderInfo struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

type GetModelsProviderInfoListResp struct {
	Providers []ModelsProviderInfo `json:"providers"`
	Total     int                  `json:"total"`
}

// UserAdminPaymentItem 充值记录项
type UserAdminPaymentItem struct {
	Id         int64  `json:"id" xorm:"id"`
	UserId     int64  `json:"user_id" xorm:"user_id"`
	UserName   string `json:"user_name" xorm:"user_name"`
	Amount     int64  `json:"amount" xorm:"pay_amount"`
	Reason     string `json:"reason" xorm:"pay_reason"`
	PayTime    int64  `json:"pay_time" xorm:"pay_time"`
	PayChannel string `json:"pay_channel" xorm:"pay_channel"`
	AdminId    int64  `json:"admin_id" xorm:"admin_user_id"`
	AdminName  string `json:"admin_name" xorm:"admin_name"`
}

type GetUserAdminPaymentListResp struct {
	List  []UserAdminPaymentItem `json:"list"`
	Total int                    `json:"total"`
}
