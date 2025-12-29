package responses

type ModelsProviderInfo struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

type GetModelsProviderInfoListResp struct {
	Providers []ModelsProviderInfo `json:"providers"`
	Total     int                  `json:"total"`
}
