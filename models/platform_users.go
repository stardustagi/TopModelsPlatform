package models

import (
	"encoding/json"
	"fmt"
	"time"
)

type PlatformUsers struct {
	Id         int64  `json:"id" xorm:"'id' pk autoincr BIGINT(12)"`
	UserName   string `json:"user_name" xorm:"'user_name' VARCHAR(64)"`
	Email      string `json:"email" xorm:"'email' VARCHAR(64)"`
	Phone      string `json:"phone" xorm:"'phone' VARCHAR(24)"`
	Password   string `json:"password" xorm:"'password' VARCHAR(128)"`
	RealName   string `json:"real_name" xorm:"'real_name' VARCHAR(128)"`
	CreatedAt  int64  `json:"created_at" xorm:"'created_at' BIGINT(12)"`
	Salt       string `json:"salt" xorm:"'salt' VARCHAR(12)"`
	LastUpdate int64  `json:"last_update" xorm:"'last_update' BIGINT(12)"`
	Active     int    `json:"active" xorm:"'active' TINYINT(1)"`
}

func (o *PlatformUsers) TableName() string {
	return "platform_users"
}

func (o *PlatformUsers) GetSliceName(slice string) string {
	return fmt.Sprintf("platform_users_%s", slice)
}

func (o *PlatformUsers) GetSliceDateMonthTable() string {
	t := time.Now()
	return fmt.Sprintf("platform_users_%d%02d", t.Year(), t.Month())
}

func (o *PlatformUsers) GetSliceDateDayTable() string {
	t := time.Now()
	return fmt.Sprintf("platform_users_%d%02d%02d", t.Year(), t.Month(), t.Day())
}

func (o *PlatformUsers) MarshalBinary() ([]byte, error) {
	return json.Marshal(o)
}

func (o *PlatformUsers) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, &o)
}

func (o *PlatformUsers) PrimaryKey() interface{} {
	return o.Id
}
