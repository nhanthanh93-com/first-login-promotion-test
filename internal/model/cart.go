package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"trinity/db/psql"
)

type Cart struct {
	Base
	UserID   uuid.UUID   `json:"user_id,omitempty" gorm:"column:user_id;index"`
	Products []*CartItem `json:"products" gorm:"foreignKey:CartID"`
}

func (m Cart) BeforeCreate(tx *gorm.DB) error {
	return nil
}

func (m Cart) TableName() string {
	return "ta_cart"
}

func (m Cart) BeforeUpdate(tx *gorm.DB) error {
	return nil
}

var DBCarts *psql.Instance[Cart]

func InitCartDB(manager *psql.DBManager) {
	DBCarts = psql.NewInstance[Cart](manager, "db1")
}
