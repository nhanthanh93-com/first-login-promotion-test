package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"trinity/db/psql"
)

type CartItem struct {
	Base
	CartID    uuid.UUID `json:"cart_id,omitempty" gorm:"column:cart_id;index"`
	ProductID uuid.UUID `json:"product_id,omitempty" gorm:"column:product_id;index"`
	Quantity  int       `json:"quantity,omitempty" gorm:"column:quantity"`
}

func (m CartItem) BeforeCreate(tx *gorm.DB) error {
	return nil
}

func (m CartItem) TableName() string {
	return "ta_cart_item"
}

func (m CartItem) BeforeUpdate(tx *gorm.DB) error {
	return nil
}

var DBCartItems *psql.Instance[CartItem]

func InitCartItemDB(manager *psql.DBManager) {
	DBCartItems = psql.NewInstance[CartItem](manager, "db1")
}
