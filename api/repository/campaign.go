package repository

import (
	"context"
	"trinity/internal/model"
	"trinity/internal/request"
	"trinity/pkg/app"
)

type CampaignRepository interface {
	Create(campaign model.Campaign) (*model.Campaign, error)
	Find(id string) (*model.Campaign, error)
	List(req *request.LCampaignReq) ([]*model.Campaign, error)
	Count(req *request.LCampaignReq) (int64, error)
	Update(campaign model.Campaign) (*model.Campaign, error)
	Delete(id string) (bool, error)
}

type campaignRepository struct {
	app *app.Config
}

func NewCampaignRepository(appConfig *app.Config) CampaignRepository {
	return &campaignRepository{app: appConfig}
}

func (r *campaignRepository) Create(campaign model.Campaign) (*model.Campaign, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.app.Timeout)
	defer cancel()

	if err := model.DBCampaigns.DB.Model(&campaign).WithContext(ctx).
		Create(&campaign).Error; err != nil {
		return nil, err
	}
	return &campaign, nil
}

func (r *campaignRepository) Find(id string) (*model.Campaign, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.app.Timeout)
	defer cancel()
	var campaign model.Campaign
	if err := model.DBCampaigns.DB.Model(&model.Campaign{}).WithContext(ctx).
		Where("id = ?", id).
		Preload("Vouchers").
		Find(&campaign).
		First(&campaign).
		Error; err != nil {
		return nil, err
	}
	return &campaign, nil
}

func (r *campaignRepository) List(req *request.LCampaignReq) ([]*model.Campaign, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.app.Timeout)
	defer cancel()
	offset := int(req.GetOffset())
	limit := int(req.GetLimit())
	order := req.GetOrder()

	var campaigns []*model.Campaign
	db := model.DBCampaigns.DB.Model(&model.Campaign{}).WithContext(ctx)

	if req.ID != "" {
		db = db.Where("id = ?", req.ID)
	}

	if req.Name != "" {
		db = db.Where("name LIKE ?", "%"+req.Name+"%")
	}

	if req.MaxUser != nil {
		db = db.Where("max_user = ?", req.MaxUser)
	}

	if req.StarDateStr != "" && req.EndDateStr != "" {
		db = db.Where("created_at BETWEEN ? AND ?", req.StartDate, req.EndDate)

	}

	if err := db.
		Offset(offset).
		Order(order).
		Limit(limit).
		Find(&campaigns).
		Error; err != nil {
		return nil, err
	}

	return campaigns, nil
}

func (r *campaignRepository) Count(req *request.LCampaignReq) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.app.Timeout)
	defer cancel()
	order := req.GetOrder()
	var count int64
	db := model.DBCampaigns.DB.Model(&model.Campaign{}).WithContext(ctx)

	if req.ID != "" {
		db = db.Where("id = ?", req.ID)
	}

	if req.Name != "" {
		db = db.Where("name LIKE ?", "%"+req.Name+"%")
	}

	if req.MaxUser != nil {
		db = db.Where("max_user = ?", req.MaxUser)
	}

	if req.StarDateStr != "" && req.EndDateStr != "" {
		db = db.Where("created_at BETWEEN ? AND ?", req.StartDate, req.EndDate)
	}

	if err := db.
		WithContext(ctx).
		Count(&count).
		Order(order).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func (r *campaignRepository) Update(campaign model.Campaign) (*model.Campaign, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.app.Timeout)
	defer cancel()

	existing, err := r.Find(campaign.ID.String())
	if err != nil {
		return nil, err
	}

	if err = model.DBCampaigns.DB.WithContext(ctx).
		Model(&model.Campaign{}).
		Where("id = ?", campaign.ID.String()).
		UpdateColumns(campaign).
		Find(&existing).
		Error; err != nil {
		return nil, err
	}

	return existing, nil
}

func (r *campaignRepository) Delete(id string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.app.Timeout)
	defer cancel()
	existing, err := r.Find(id)
	if err != nil {
		return false, err
	}

	if err = model.DBCampaigns.DB.Model(&model.Campaign{}).
		WithContext(ctx).
		Where("id = ?", id).
		Delete(&existing).
		Error; err != nil {
		return false, err
	}
	return true, nil
}
