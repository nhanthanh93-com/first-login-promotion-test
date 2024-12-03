package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Base struct {
	ID        uuid.UUID `json:"id" gorm:"column:id;primaryKey;unique;default:uuid_generate_v4()"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at;default:current_timestamp"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at;default:current_timestamp"`
}

func (m *Base) BeforeCreate(tx *gorm.DB) (err error) {
	m.ID = uuid.New()
	m.CreatedAt = time.Now()
	return
}

type Paging struct {
	Page  int   `json:"page" form:"page"`
	Limit int   `json:"limit" form:"limit"`
	Total int64 `json:"total" form:"-"`
}

func (p *Paging) Process() {
	if p.Page <= 0 {
		p.Page = 1
	}
	if p.Limit <= 0 || p.Limit >= 100 {
		p.Limit = 5
	}
}

type JSONB []interface{}
