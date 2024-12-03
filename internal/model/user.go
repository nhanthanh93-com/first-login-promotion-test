package model

import (
	"gorm.io/gorm"
	"time"
	"trinity/db/psql"
)

type User struct {
	Base
	Email      string     `json:"email,omitempty" gorm:"column:email;uniqueIndex:idx_email,unique" validate:"required"`
	RegisterAt time.Time  `json:"register_at,omitempty" gorm:"column:register_at;default:current_timestamp"`
	Vouchers   []*Voucher `json:"vouchers,omitempty" gorm:"foreignKey:UserID"`
}

func (m User) BeforeCreate(tx *gorm.DB) error {
	return nil
}

func (m User) TableName() string {
	return "ta_users"
}

func (m User) BeforeUpdate(tx *gorm.DB) error {
	return nil
}

var DBUsers *psql.Instance[User]

func InitUserDB(manager *psql.DBManager) {
	DBUsers = psql.NewInstance[User](manager, "db1")
}
