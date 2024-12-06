package repository

import (
	"context"
	"trinity/internal/model"
	"trinity/internal/request"
	"trinity/pkg/app"
)

type OrderRepository interface {
	Create(order model.Order) (*model.Order, error)
	Find(id string) (*model.Order, error)
	UpdateOrderStatus(req *request.UOrderStatusReq) (*model.Order, error)
}

type orderRepository struct {
	app *app.Config
}

func NewOrderRepository(appConfig *app.Config) OrderRepository {
	return &orderRepository{appConfig}
}

func (r *orderRepository) Create(order model.Order) (*model.Order, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.app.Timeout)
	defer cancel()
	if err := model.DBOrders.DB.Model(&order).WithContext(ctx).
		Create(&order).Error; err != nil {
		return nil, err
	}

	return &order, nil
}

func (r *orderRepository) Find(id string) (*model.Order, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.app.Timeout)
	defer cancel()
	var order model.Order
	if err := model.DBOrders.DB.Model(&model.Order{}).WithContext(ctx).
		Where("id = ?", id).
		Find(&order).
		First(&order).
		Error; err != nil {
		return nil, err
	}

	return &order, nil
}

func (r *orderRepository) UpdateOrderStatus(req *request.UOrderStatusReq) (*model.Order, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.app.Timeout)
	defer cancel()

	existing, err := r.Find(req.ID)
	if err != nil {
		return nil, err
	}

	if err = model.DBOrders.DB.WithContext(ctx).
		Model(&model.Order{}).
		Where("id = ?", req.ID).
		Update("status", req.Status).
		Find(&existing).
		Error; err != nil {
		return nil, err
	}

	return existing, nil
}
