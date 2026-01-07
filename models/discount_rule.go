package models

import (
	"encoding/json"
	"fmt"
	"time"
)

type DiscountRule struct {
	Id           int64  `json:"id" xorm:"'id' pk autoincr BIGINT(20)"`
	Name         string `json:"name" xorm:"'name' not null default '' VARCHAR(128)"`
	Description  string `json:"description" xorm:"'description' not null default '' VARCHAR(512)"`
	DiscountRate int    `json:"discount_rate" xorm:"'discount_rate' not null default 0 INT(10)"`
	Status       int    `json:"status" xorm:"'status' not null default 1 INT(10)"`
	CreatedAt    int64  `json:"created_at" xorm:"'created_at' not null default 0 BIGINT(20)"`
	UpdatedAt    int64  `json:"updated_at" xorm:"'updated_at' not null default 0 BIGINT(20)"`
}

func (o *DiscountRule) TableName() string {
	return "discount_rule"
}

func (o *DiscountRule) GetSliceName(slice string, num uint32) string {
	var hash uint32
	for _, c := range slice {
		hash = hash*31 + uint32(c)
	}
	shardIndex := hash % num
	return fmt.Sprintf("discount_rule_%d", shardIndex)
}

func (o *DiscountRule) GetSliceDateMonthTable() string {
	t := time.Now()
	return fmt.Sprintf("discount_rule_%d%02d", t.Year(), t.Month())
}

func (o *DiscountRule) GetSliceDateDayTable() string {
	t := time.Now()
	return fmt.Sprintf("discount_rule_%d%02d%02d", t.Year(), t.Month(), t.Day())
}

func (o *DiscountRule) MarshalBinary() ([]byte, error) {
	return json.Marshal(o)
}

func (o *DiscountRule) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, &o)
}

func (o *DiscountRule) PrimaryKey() interface{} {
	return o.Id
}
