package repository

import (
	"context"
	"trinity/internal/model"
	"trinity/pkg/app"
)

type CartRepository interface {
	Create(cart model.Cart) (*model.Cart, error)
	Find(id string) (*model.Cart, error)
}

type cartRepository struct {
	app *app.Config
}

func NewCartRepository(appConfig *app.Config) CartRepository {
	return &cartRepository{appConfig}
}

func (r *cartRepository) Create(cart model.Cart) (*model.Cart, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.app.Timeout)
	defer cancel()

	if err := model.DBCarts.DB.Model(&cart).WithContext(ctx).
		Create(&cart).Error; err != nil {
		return nil, err
	}

	return &cart, nil
}

func (r *cartRepository) Find(id string) (*model.Cart, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.app.Timeout)
	defer cancel()
	var item model.Cart
	if err := model.DBCarts.DB.Model(&model.Cart{}).WithContext(ctx).
		Preload("Products").
		Where("id = ?", id).
		Find(&item).
		First(&item).
		Error; err != nil {
		return nil, err
	}

	return &item, nil
}
