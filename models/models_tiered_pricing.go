package models

import (
	"encoding/json"
	"fmt"
	"time"
)

type ModelsTieredPricing struct {
	Id          int64 `json:"id" xorm:"'id' pk autoincr BIGINT(20)"`
	ModelId     int64 `json:"model_id" xorm:"'model_id' BIGINT(20)"`
	TierStart   int64 `json:"tier_start" xorm:"'tier_start' BIGINT(20)"`
	TierEnd     int64 `json:"tier_end" xorm:"'tier_end' BIGINT(20)"`
	InputPrice  int   `json:"input_price" xorm:"'input_price' INT(10)"`
	OutputPrice int   `json:"output_price" xorm:"'output_price' INT(10)"`
	CachePrice  int   `json:"cache_price" xorm:"'cache_price' INT(10)"`
	CreatedAt   int64 `json:"created_at" xorm:"'created_at' BIGINT(20)"`
	UpdatedAt   int64 `json:"updated_at" xorm:"'updated_at' BIGINT(20)"`
}

func (o *ModelsTieredPricing) TableName() string {
	return "models_tiered_pricing"
}

func (o *ModelsTieredPricing) GetSliceName(slice string, num uint32) string {
	var hash uint32
	for _, c := range slice {
		hash = hash*31 + uint32(c)
	}
	shardIndex := hash % num
	return fmt.Sprintf("models_tiered_pricing_%d", shardIndex)
}

func (o *ModelsTieredPricing) GetSliceDateMonthTable() string {
	t := time.Now()
	return fmt.Sprintf("models_tiered_pricing_%d%02d", t.Year(), t.Month())
}

func (o *ModelsTieredPricing) GetSliceDateDayTable() string {
	t := time.Now()
	return fmt.Sprintf("models_tiered_pricing_%d%02d%02d", t.Year(), t.Month(), t.Day())
}

func (o *ModelsTieredPricing) MarshalBinary() ([]byte, error) {
	return json.Marshal(o)
}

func (o *ModelsTieredPricing) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, &o)
}

func (o *ModelsTieredPricing) PrimaryKey() interface{} {
	return o.Id
}
