package repository

import (
	"context"
	"trinity/internal/model"
	"trinity/internal/request"
	"trinity/pkg/app"
)

type ProductRepository interface {
	Create(product model.Product) (*model.Product, error)
	Find(id string) (*model.Product, error)
	List(req *request.LProductReq) ([]*model.Product, error)
	Count(req *request.LProductReq) (int64, error)
	Update(product model.Product) (*model.Product, error)
	Delete(id string) (bool, error)
}

type productRepository struct {
	app *app.Config
}

func NewProductRepository(appConfig *app.Config) ProductRepository {
	return &productRepository{app: appConfig}
}

func (r *productRepository) Create(product model.Product) (*model.Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.app.Timeout)
	defer cancel()

	if err := model.DBProducts.DB.Model(&product).WithContext(ctx).
		Create(&product).Error; err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *productRepository) Find(id string) (*model.Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.app.Timeout)
	defer cancel()
	var product model.Product
	if err := model.DBProducts.DB.Model(&model.Product{}).WithContext(ctx).
		Where("id = ?", id).
		Find(&product).
		First(&product).
		Error; err != nil {
		return nil, err
	}

	return &product, nil
}

func (r *productRepository) List(req *request.LProductReq) ([]*model.Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.app.Timeout)
	defer cancel()

	offset := int(req.GetOffset())
	limit := int(req.GetLimit())
	order := req.GetOrder()

	var products []*model.Product
	db := model.DBProducts.DB.Model(&model.Product{}).WithContext(ctx)

	if req.Name != "" {
		db = db.Where("name = ?", req.Name)
	}
	if req.ID != "" {
		db = db.Where("id = ?", req.ID)
	}

	if req.StarDateStr != "" && req.EndDateStr != "" {
		db = db.Where("created_at BETWEEN ? AND ?", req.StartDate, req.EndDate)
	}

	if err := db.
		Offset(offset).
		Order(order).
		Limit(limit).
		Find(&products).
		Error; err != nil {
		return nil, err
	}
	return products, nil
}

func (r *productRepository) Count(req *request.LProductReq) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.app.Timeout)
	defer cancel()
	order := req.GetOrder()
	var count int64
	db := model.DBProducts.DB.Model(&model.Product{}).WithContext(ctx)

	if req.Name != "" {
		db = db.Where("name = ?", req.Name)
	}
	if req.ID != "" {
		db = db.Where("id = ?", req.ID)
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

func (r *productRepository) Update(product model.Product) (*model.Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.app.Timeout)
	defer cancel()

	existing, err := r.Find(product.ID.String())
	if err != nil {
		return nil, err
	}

	if err = model.DBProducts.DB.WithContext(ctx).
		Model(&model.Product{}).
		Where("id = ?", product.ID.String()).
		UpdateColumns(&product).
		Find(&existing).
		Error; err != nil {
		return nil, err
	}

	return existing, nil
}

func (r *productRepository) Delete(id string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.app.Timeout)
	defer cancel()
	existing, err := r.Find(id)
	if err != nil {
		return false, err
	}

	if err = model.DBProducts.DB.Model(&model.Product{}).
		WithContext(ctx).
		Where("id = ?", id).
		Delete(&existing).
		Error; err != nil {
		return false, err
	}
	return true, nil
}
