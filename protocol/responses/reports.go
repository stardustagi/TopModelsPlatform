package responses

import "time"

type DayReportItem struct {
	Id               int64  `json:"id"`
	UserId           int64  `json:"user_id"`
	UserName         string `json:"user_name"`
	ActualProviderId int    `json:"actual_provider_id"`
	ProviderName     string `json:"provider_name"`
	ModelId          int    `json:"model_id"`
	ModelName        string `json:"model_name"`
	ConsumeType      string `json:"consume_type"`
	Date             string `json:"date"`
	TotalConsumed    int64  `json:"total_consumed"`
	TotalCost        int64  `json:"total_cost"`
	UpdatedAt        int64  `json:"updated_at"`
}

type GetDayReportListResp struct {
	Items []DayReportItem `json:"items"`
	Total int             `json:"total"`
	Skip  int             `json:"skip"`
	Limit int             `json:"limit"`
}

// UserCallStatusItem 用户调用状态统计项
type UserCallStatusItem struct {
	UserId             int64   `json:"user_id" xorm:"user_id"`
	CallerKey          string  `json:"caller_key" xorm:"caller_key"`
	Model              string  `json:"model" xorm:"model"`
	CallCount5Min      int64   `json:"call_count_5min" xorm:"call_count_5min"`
	SuccessCount5Min   int64   `json:"success_count_5min" xorm:"success_count_5min"`
	SuccessRate5Min    float64 `json:"success_rate_5min" xorm:"success_rate_5min"`
	AvgLatency5Min     float64 `json:"avg_latency_5min" xorm:"avg_latency_5min"`
	CallCount10Min     int64   `json:"call_count_10min" xorm:"call_count_10min"`
	SuccessCount10Min  int64   `json:"success_count_10min" xorm:"success_count_10min"`
	SuccessRate10Min   float64 `json:"success_rate_10min" xorm:"success_rate_10min"`
	AvgLatency10Min    float64 `json:"avg_latency_10min" xorm:"avg_latency_10min"`
	CallCount30Min     int64   `json:"call_count_30min" xorm:"call_count_30min"`
	SuccessCount30Min  int64   `json:"success_count_30min" xorm:"success_count_30min"`
	SuccessRate30Min   float64 `json:"success_rate_30min" xorm:"success_rate_30min"`
	AvgLatency30Min    float64 `json:"avg_latency_30min" xorm:"avg_latency_30min"`
	CallCount1Hour     int64   `json:"call_count_1hour" xorm:"call_count_1hour"`
	SuccessCount1Hour  int64   `json:"success_count_1hour" xorm:"success_count_1hour"`
	SuccessRate1Hour   float64 `json:"success_rate_1hour" xorm:"success_rate_1hour"`
	AvgLatency1Hour    float64 `json:"avg_latency_1hour" xorm:"avg_latency_1hour"`
	CallCount24Hour    int64   `json:"call_count_24hour" xorm:"call_count_24hour"`
	SuccessCount24Hour int64   `json:"success_count_24hour" xorm:"success_count_24hour"`
	SuccessRate24Hour  float64 `json:"success_rate_24hour" xorm:"success_rate_24hour"`
	AvgLatency24Hour   float64 `json:"avg_latency_24hour" xorm:"avg_latency_24hour"`
}

// GetDayReportSummaryResp 调用统计汇总响应
type GetDayReportSummaryResp struct {
	List  []UserCallStatusItem `json:"list"`
	Total int                  `json:"total"`
}

type GetReportDetailListItem struct {
	Id               uint64    `json:"id" xorm:"'id' not null pk autoincr comment('主键ID') UNSIGNED BIGINT(20)"`
	Model            string    `json:"model" xorm:"'model' not null default '' comment('模型名字') index VARCHAR(64)"`
	ModelId          int       `json:"model_id" xorm:"'model_id' not null default 0 comment('模型ID（计费使用）') INT(11)"`
	ActualModel      string    `json:"actual_model" xorm:"'actual_model' not null default '' comment('实际使用的模型') VARCHAR(64)"`
	Provider         string    `json:"provider" xorm:"'provider' not null default '' comment('虚拟provider') VARCHAR(64)"`
	ActualProvider   string    `json:"actual_provider" xorm:"'actual_provider' not null default '' comment('实际服务商') VARCHAR(64)"`
	ActualProviderId string    `json:"actual_provider_id" xorm:"'actual_provider_id' not null default '' comment('实际服务商ID') VARCHAR(64)"`
	CallerKey        string    `json:"caller_key" xorm:"'caller_key' not null default '' comment('客户端key') index VARCHAR(128)"`
	Stream           int       `json:"stream" xorm:"'stream' not null default 0 comment('是否流式访问：0-否，1-是') TINYINT(1)"`
	ReportType       string    `json:"report_type" xorm:"'report_type' not null default '' comment('报告类型：text/image/video') VARCHAR(16)"`
	TokensPerSec     int       `json:"tokens_per_sec" xorm:"'tokens_per_sec' not null default 0 comment('每秒输出token') INT(11)"`
	Latency          string    `json:"latency" xorm:"'latency' not null default 0.0000 comment('请求延迟（秒）') DECIMAL(10,4)"`
	Step             string    `json:"step" xorm:"'step' not null default '' comment('调用环节：call_llm_agent/check_user_balance/select_provider/send_llm_request/send_llm_completed/llm_agent_done/user_agent_done') index VARCHAR(32)"`
	StatusCode       string    `json:"status_code" xorm:"'status_code' not null default '' comment('状态码（非空为失败）') VARCHAR(16)"`
	StatusMessage    string    `json:"status_message" xorm:"'status_message' not null default '' comment('状态消息（状态码非空时有值）') VARCHAR(512)"`
	CreatedAt        time.Time `json:"created_at" xorm:"'created_at' not null default current_timestamp() comment('请求时间') index DATETIME"`
	UserId           int64     `json:"user_id" xorm:"'user_id' BIGINT(20)"`
	UserName         string    `json:"user_name" xorm:"'user_name' VARCHAR(64)"`
}
type GetReportDetailListResp struct {
	Items []GetReportDetailListItem `json:"items"`
	Total int                       `json:"total"`
}
