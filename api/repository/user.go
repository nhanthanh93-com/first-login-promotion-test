package repository

import (
	"context"
	"fmt"
	"strings"
	"trinity/helpers"
	"trinity/internal/model"
	"trinity/internal/request"
	"trinity/pkg/app"
)

type UserRepository interface {
	Register(req *request.RUserReq) (*model.User, error)
	Create(m model.User) (*model.User, error)
	Find(id string) (*model.User, error)
	List(r *request.LUserReq) ([]*model.User, error)
	Count(req *request.LUserReq) (int64, error)
	Update(user model.User) (*model.User, error)
	Delete(id string) (bool, error)
}

type userRepository struct {
	app *app.Config
}

func NewUserRepository(appConfig *app.Config) UserRepository {
	return &userRepository{app: appConfig}
}

func (r *userRepository) Register(req *request.RUserReq) (*model.User, error) {
	var err error
	var campaign model.Campaign
	var totalVouchers int64
	var voucher model.Voucher
	wg := helpers.NewWgGroup()
	wg.Go(func() error {
		ctx, cancel := context.WithTimeout(context.Background(), r.app.Timeout)
		defer cancel()
		if err = model.DBCampaigns.DB.Model(&campaign).WithContext(ctx).
			Where("id = ?", req.CampaignID).
			Find(&campaign).
			First(&campaign).
			Error; err != nil {
			return err
		}
		return nil
	})

	wg.Go(func() error {
		ctx, cancel := context.WithTimeout(context.Background(), r.app.Timeout)
		defer cancel()
		if err = model.DBVouchers.DB.Model(voucher).WithContext(ctx).
			Where("campaign_id = ?", req.CampaignID).
			Count(&totalVouchers).Error; err != nil {
			return err
		}
		return nil
	})
	if err = wg.Wait(); err != nil {
		return nil, nil
	}

	if int32(totalVouchers) >= campaign.MaxUser {
		return nil, fmt.Errorf("error: Registration limit exceeded")
	}

	var m model.User
	m.Email = req.Email

	user, err := r.Create(m)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), r.app.Timeout)
	defer cancel()

	voucher.CampaignID = campaign.ID
	voucher.UserID = user.ID
	voucher.ExpiresAt = campaign.ExpiresAt
	voucher.Discount = campaign.Discount
	if err = model.DBVouchers.DB.Model(&voucher).WithContext(ctx).
		Create(&voucher).Error; err != nil {
		return nil, err
	}
	voucher.Campaign = nil
	voucher.User = nil

	user.Vouchers = append(user.Vouchers, &voucher)
	return user, nil
}

func (r *userRepository) Create(user model.User) (*model.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.app.Timeout)
	defer cancel()

	if err := model.DBUsers.DB.Model(&user).WithContext(ctx).
		Create(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Find(id string) (*model.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.app.Timeout)
	defer cancel()
	var user model.User
	if err := model.DBUsers.DB.Model(&model.User{}).WithContext(ctx).
		Where("id = ?", id).
		Find(&user).
		First(&user).
		Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) List(req *request.LUserReq) ([]*model.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.app.Timeout)
	defer cancel()

	offset := int(req.GetOffset())
	limit := int(req.GetLimit())
	order := req.GetOrder()

	var users []*model.User
	db := model.DBUsers.DB.Model(&model.User{}).WithContext(ctx)

	if req.Email != "" {
		db = db.Where("email = ?", req.Email)
	}
	if req.ID != "" {
		db = db.Where("id = ?", req.ID)
	}

	if req.StarDateStr != "" && req.EndDateStr != "" {
		db = db.Where("created_at BETWEEN ? AND ?", req.StartDate, req.EndDate)
	}

	if req.Vouchers != "" {
		codes := strings.Split(strings.ReplaceAll(req.Vouchers, " ", ""), ",")
		db = db.Joins("JOIN ta_voucher ON ta_voucher.user_id = ta_users.id").
			Where("ta_voucher.code IN ?", codes).Group("ta_users.id")
	}

	if err := db.
		Preload("Vouchers").
		Offset(offset).
		Order(order).
		Limit(limit).
		Find(&users).
		Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (r *userRepository) Count(req *request.LUserReq) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.app.Timeout)
	defer cancel()
	order := req.GetOrder()
	var count int64
	db := model.DBUsers.DB.Model(&model.User{})

	if req.Email != "" {
		db = db.Where("email = ?", req.Email)
	}
	if req.ID != "" {
		db = db.Where("id = ?", req.ID)
	}

	if req.StarDateStr != "" && req.EndDateStr != "" {
		db = db.Where("created_at BETWEEN ? AND ?", req.StartDate, req.EndDate)
	}

	if req.Vouchers != "" {
		codes := strings.Split(strings.ReplaceAll(req.Vouchers, " ", ""), ",")
		db = db.Joins("JOIN ta_voucher ON ta_voucher.user_id = ta_users.id").
			Where("ta_voucher.code IN ?", codes).Group("ta_users.id")
	}

	if err := db.
		WithContext(ctx).
		Count(&count).
		Order(order).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func (r *userRepository) Update(user model.User) (*model.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.app.Timeout)
	defer cancel()

	existingUser, err := r.Find(user.ID.String())
	if err != nil {
		return nil, err
	}

	if err = model.DBUsers.DB.WithContext(ctx).
		Model(&model.User{}).
		Where("id = ?", user.ID.String()).
		UpdateColumns(user).
		Find(&existingUser).
		Error; err != nil {
		return nil, err
	}

	return existingUser, nil
}

func (r *userRepository) Delete(id string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.app.Timeout)
	defer cancel()
	existingUser, err := r.Find(id)
	if err != nil {
		return false, fmt.Errorf("user not found: %w", err)
	}

	if err = model.DBUsers.DB.Model(&model.User{}).
		WithContext(ctx).
		Where("id = ?", id).
		Delete(&existingUser).
		Error; err != nil {
		return false, fmt.Errorf("failed to delete user: %w", err)
	}
	return true, nil
}
