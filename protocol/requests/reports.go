package requests

type GetDayReportSummaryReq struct {
	PageInfo PageReq `json:"page_info"`
	UserId   int64   `json:"user_id"`
}

type GetReportDetailListReq struct {
	PageInfo  PageReq `json:"page_info"`
	StartDate int64   `json:"start_date" validate:"required"`
	EndDate   int64   `json:"end_date" validate:"required"`
	CallKey   string  `json:"call_key"`
	Model     string  `json:"model"`
}
