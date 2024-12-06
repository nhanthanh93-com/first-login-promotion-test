package model

import (
	"gorm.io/gorm"
	"trinity/db/psql"
)

type Product struct {
	Base
	Name      string      `json:"name,omitempty" gorm:"column:name;" validate:"required"`
	Price     float64     `json:"price,omitempty" gorm:"column:price;default:0"`
	Stock     int         `json:"stock,omitempty" gorm:"column:stock;default:0"`
	Campaigns []*Campaign `json:"campaigns,omitempty" gorm:"many2many:ta_campaign_product;"`
}

func (m Product) BeforeCreate(tx *gorm.DB) error {
	return nil
}

func (m Product) TableName() string {
	return "ta_product"
}

func (m Product) BeforeUpdate(tx *gorm.DB) error {
	return nil
}

var DBProducts *psql.Instance[Product]

func InitProductDB(manager *psql.DBManager) {
	DBProducts = psql.NewInstance[Product](manager, "db1")
}
