package models

import (
	"encoding/json"
	"fmt"
	"time"
)

type UserDiscount struct {
	Id           int64 `json:"id" xorm:"'id' pk autoincr BIGINT(20)"`
	UserId       int64 `json:"user_id" xorm:"'user_id' not null default 0 BIGINT(20) unique(uk_user_model_rule)" `
	ModelId      int   `json:"model_id" xorm:"'model_id' not null default 0 INT(10) unique(uk_user_model_rule)"`
	RuleId       int64 `json:"rule_id" xorm:"'rule_id' not null default 0 BIGINT(20) unique(uk_user_model_rule)"`
	DiscountRate int   `json:"discount_rate" xorm:"'discount_rate' not null default 0 INT(10)"`
	CreatedAt    int64 `json:"created_at" xorm:"'created_at' not null default 0 BIGINT(20)"`
	UpdatedAt    int64 `json:"updated_at" xorm:"'updated_at' not null default 0 BIGINT(20)"`
}

func (o *UserDiscount) TableName() string {
	return "user_discount"
}

func (o *UserDiscount) GetSliceName(slice string, num uint32) string {
	var hash uint32
	for _, c := range slice {
		hash = hash*31 + uint32(c)
	}
	shardIndex := hash % num
	return fmt.Sprintf("user_discount_%d", shardIndex)
}

func (o *UserDiscount) GetSliceDateMonthTable() string {
	t := time.Now()
	return fmt.Sprintf("user_discount_%d%02d", t.Year(), t.Month())
}

func (o *UserDiscount) GetSliceDateDayTable() string {
	t := time.Now()
	return fmt.Sprintf("user_discount_%d%02d%02d", t.Year(), t.Month(), t.Day())
}

func (o *UserDiscount) MarshalBinary() ([]byte, error) {
	return json.Marshal(o)
}

func (o *UserDiscount) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, &o)
}

func (o *UserDiscount) PrimaryKey() interface{} {
	return o.Id
}
