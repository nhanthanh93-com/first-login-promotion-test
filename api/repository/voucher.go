package repository

import (
	"context"
	"fmt"
	"time"
	"trinity/internal/model"
	"trinity/internal/request"
	"trinity/pkg/app"
)

type VoucherRepository interface {
	Create(voucher model.Voucher) (*model.Voucher, error)
	Find(id string) (*model.Voucher, error)
	List(req *request.LVoucherReq) ([]*model.Voucher, error)
	Count(req *request.LVoucherReq) (int64, error)
	Update(voucher model.Voucher) (*model.Voucher, error)
	Delete(id string) (bool, error)
}

type voucherRepository struct {
	app *app.Config
}

func NewVoucherRepository(appConfig *app.Config) VoucherRepository {
	return &voucherRepository{app: appConfig}
}

func (r *voucherRepository) Create(Voucher model.Voucher) (*model.Voucher, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.app.Timeout)
	defer cancel()

	if err := model.DBVouchers.DB.Model(&Voucher).WithContext(ctx).
		Create(&Voucher).Error; err != nil {
		return nil, err
	}
	return &Voucher, nil
}

func (r *voucherRepository) Find(id string) (*model.Voucher, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.app.Timeout)
	defer cancel()
	var voucher model.Voucher
	if err := model.DBVouchers.DB.Model(&model.Voucher{}).WithContext(ctx).
		Where("id = ?", id).
		Find(&voucher).
		First(&voucher).
		Error; err != nil {
		return nil, err
	}
	return &voucher, nil
}

func (r *voucherRepository) List(req *request.LVoucherReq) ([]*model.Voucher, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.app.Timeout)
	defer cancel()
	offset := int(req.GetOffset())
	limit := int(req.GetLimit())
	order := req.GetOrder()

	var Vouchers []*model.Voucher
	db := model.DBVouchers.DB.Model(&model.Voucher{}).WithContext(ctx)

	if req.ID != "" {
		db = db.Where("id = ?", req.ID)
	}

	if req.Code != "" {
		db = db.Where("code = ?", req.Code)
	}

	if req.StarDateStr != "" && req.EndDateStr != "" {
		db = db.Where("created_at BETWEEN ? AND ?", req.StartDate, req.EndDate)
	}

	if err := db.
		Preload("User").
		Preload("Campaign").
		Offset(offset).
		Order(order).
		Limit(limit).
		Find(&Vouchers).
		Error; err != nil {
		return nil, err
	}

	return Vouchers, nil
}

func (r *voucherRepository) Count(req *request.LVoucherReq) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.app.Timeout)
	defer cancel()
	order := req.GetOrder()
	var count int64
	db := model.DBVouchers.DB.Model(&model.Voucher{}).WithContext(ctx)

	if req.ID != "" {
		db = db.Where("id = ?", req.ID)
	}

	if req.Code != "" {
		db = db.Where("code = ?", req.Code)
	}

	if req.StarDateStr != "" && req.EndDateStr != "" {
		db = db.Where("created_at BETWEEN ? AND ?", req.StartDate, req.EndDate)
	}

	if err := db.
		Count(&count).
		Order(order).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func (r *voucherRepository) Update(Voucher model.Voucher) (*model.Voucher, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.app.Timeout)
	defer cancel()

	existing, err := r.Find(Voucher.ID.String())
	if err != nil {
		return nil, err
	}
	currentTime := time.Now()
	if currentTime.After(existing.ExpiresAt) {
		return nil, fmt.Errorf("error: Your voucher (Code: [%v]) has expired and can no longer be redeemed", existing.Code)
	}

	if err = model.DBVouchers.DB.WithContext(ctx).
		Model(&model.Voucher{}).
		Where("id = ?", Voucher.ID.String()).
		UpdateColumns(Voucher).
		Find(&existing).
		Error; err != nil {
		return nil, err
	}

	return existing, nil
}

func (r *voucherRepository) Delete(id string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.app.Timeout)
	defer cancel()
	existing, err := r.Find(id)
	if err != nil {
		return false, err
	}

	if err = model.DBVouchers.DB.Model(&model.Voucher{}).
		WithContext(ctx).
		Where("id = ?", id).
		Delete(&existing).
		Error; err != nil {
		return false, err
	}
	return true, nil
}
