package model

import (
	"gorm.io/gorm"
	"time"
	"trinity/db/psql"
)

type Campaign struct {
	Base
	Name      string     `json:"name,omitempty" gorm:"column:name;" validate:"required"`
	MaxUser   int32      `json:"max_user,omitempty" gorm:"column:max_user;" validate:"required"`
	ExpiresAt time.Time  `json:"expires_at,omitempty" gorm:"column:expires_at;default:current_timestamp"`
	Discount  int32      `json:"discount,omitempty" gorm:"column:discount;default:0"`
	Vouchers  []*Voucher `json:"vouchers,omitempty" gorm:"foreignKey:CampaignID"`
	Products  []*Product `json:"products,omitempty" gorm:"many2many:ta_campaign_product;"`
}

func (m Campaign) BeforeCreate(tx *gorm.DB) error {
	return nil
}

func (m Campaign) TableName() string {
	return "ta_campaign"
}

func (m Campaign) BeforeUpdate(tx *gorm.DB) error {
	return nil
}

var DBCampaigns *psql.Instance[Campaign]

func InitCampaignDB(manager *psql.DBManager) {
	DBCampaigns = psql.NewInstance[Campaign](manager, "db1")
}
