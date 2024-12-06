package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"trinity/db/psql"
)

type Order struct {
	Base
	UserID    uuid.UUID `json:"user_id,omitempty" gorm:"column:user_id;index"`
	Total     float64   `json:"total"`
	Status    string    `json:"status"`
	CreatedAt int64     `json:"created_at"`
	UpdatedAt int64     `json:"updated_at"`
}

func (m Order) BeforeCreate(tx *gorm.DB) error {
	return nil
}

func (m Order) TableName() string {
	return "ta_order"
}

func (m Order) BeforeUpdate(tx *gorm.DB) error {
	return nil
}

var DBOrders *psql.Instance[Order]

func InitOrderDB(manager *psql.DBManager) {
	DBOrders = psql.NewInstance[Order](manager, "db1")
}
