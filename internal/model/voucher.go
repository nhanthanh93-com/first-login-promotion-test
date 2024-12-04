package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
	"trinity/db/psql"
	"trinity/pkg/nanoid"
)

type Voucher struct {
	Base
	Code       string    `json:"code,omitempty" gorm:"column:code;uniqueIndex:idx_code,unique" validate:"required"`
	ExpiresAt  time.Time `json:"expires_at,omitempty" gorm:"column:expires_at;default:current_timestamp"`
	IsUsed     *bool     `json:"is_used,omitempty" gorm:"column:is_used;default:false"`
	Discount   int32     `json:"discount,omitempty" gorm:"column:discount;default:0"`
	CampaignID uuid.UUID `json:"-,omitempty" gorm:"column:campaign_id;index"`
	Campaign   *Campaign `json:"campaign,omitempty" gorm:"foreignKey:CampaignID"`
	UserID     uuid.UUID `json:"-,omitempty" gorm:"column:user_id;index"`
	User       *User     `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

func (m Voucher) BeforeCreate(tx *gorm.DB) error {
	code, err := nanoid.New()
	if err != nil {
		return err
	}

	tx.Statement.SetColumn("code", code)
	return nil
}

func (m Voucher) TableName() string {
	return "ta_voucher"
}

func (m Voucher) BeforeUpdate(tx *gorm.DB) error {
	return nil
}

var DBVouchers *psql.Instance[Voucher]

func InitVoucherDB(manager *psql.DBManager) {
	DBVouchers = psql.NewInstance[Voucher](manager, "db1")
}
