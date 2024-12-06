package model

import (
	"gorm.io/gorm"
	"trinity/db/psql"
)

type CampaignsProducts struct {
	CampaignID string `json:"campaign_id" gorm:"column:campaign_id;primaryKey;"`
	ProductID  string `json:"product_id" gorm:"column:product_id;primaryKey;"`
}

func (CampaignsProducts) BeforeCreate(db *gorm.DB) error {
	return nil
}

func (CampaignsProducts) TableName() string {
	return "ta_campaign_product"
}

var DBCampaignsProducts *psql.Instance[CampaignsProducts]

func InitCampaignsProductsDB(manager *psql.DBManager) {
	DBCampaignsProducts = psql.NewInstance[CampaignsProducts](manager, "db1")
}
