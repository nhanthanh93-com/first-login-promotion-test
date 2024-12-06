package repository

import (
	"context"
	"trinity/internal/model"
	"trinity/pkg/app"
)

type CartItemRepository interface {
	AddToCart(item model.CartItem) (*model.CartItem, error)
	Find(id string) (*model.CartItem, error)
	Update(item model.CartItem) (*model.CartItem, error)
	Delete(itemID string) (bool, error)
	Deletes(products []*model.CartItem) (bool, error)
}

type cartItemRepository struct {
	app *app.Config
}

func NewCartItemRepository(appConfig *app.Config) CartItemRepository {
	return &cartItemRepository{appConfig}
}

func (r *cartItemRepository) AddToCart(item model.CartItem) (*model.CartItem, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.app.Timeout)
	defer cancel()

	if err := model.DBCartItems.DB.Model(&item).WithContext(ctx).
		Create(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *cartItemRepository) Find(id string) (*model.CartItem, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.app.Timeout)
	defer cancel()
	var item model.CartItem
	if err := model.DBCartItems.DB.Model(&model.CartItem{}).WithContext(ctx).
		Where("id = ?", id).
		Find(&item).
		First(&item).
		Error; err != nil {
		return nil, err
	}

	return &item, nil
}

func (r *cartItemRepository) Update(item model.CartItem) (*model.CartItem, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.app.Timeout)
	defer cancel()
	existing, err := r.Find(item.ID.String())
	if err != nil {
		return nil, err
	}

	if err = model.DBCartItems.DB.Model(&item).WithContext(ctx).
		Where("id = ?", item.ID.String()).
		UpdateColumns(&item).
		Find(&existing).
		Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *cartItemRepository) Delete(itemID string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.app.Timeout)
	defer cancel()
	existing, err := r.Find(itemID)
	if err != nil {
		return false, err
	}

	if err = model.DBProducts.DB.Model(&model.Product{}).
		WithContext(ctx).
		Where("id = ?", itemID).
		Delete(&existing).
		Error; err != nil {
		return false, err
	}
	return true, nil
}

func (r *cartItemRepository) Deletes(products []*model.CartItem) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.app.Timeout)
	defer cancel()

	if err := model.DBCartItems.DB.Model(&model.CartItem{}).
		WithContext(ctx).
		Delete(&products).
		Error; err != nil {
		return false, err
	}
	return true, nil
}
