package requests

type ModelIdReq struct {
	Id int64 `json:"id" query:"id" validate:"required"`
}

type CreateModelReq struct {
	Name        string `json:"name" validate:"required"`
	ApiVersion  string `json:"api_version"`
	DeployName  string `json:"deploy_name"`
	InputPrice  int    `json:"input_price"`
	OutputPrice int    `json:"output_price"`
	CachePrice  int    `json:"cache_price"`
	IsPrivate   int    `json:"is_private"`
	OwnerId     int64  `json:"owner_id"`
	Address     string `json:"address"`
	ApiStyles   string `json:"api_styles"`
}

type UpdateModelReq struct {
	Id          int64  `json:"id" validate:"required"`
	Name        string `json:"name"`
	ApiVersion  string `json:"api_version"`
	DeployName  string `json:"deploy_name"`
	InputPrice  int    `json:"input_price"`
	OutputPrice int    `json:"output_price"`
	Status      string `json:"status"`
	Address     string `json:"address"`
	ApiStyles   string `json:"api_styles"`
}
