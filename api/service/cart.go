package service

import (
	"trinity/api/repository"
	"trinity/internal/model"
	"trinity/internal/request"
	"trinity/pkg/app"
)

type CartService interface {
	AddToCart(item model.CartItem) (*model.CartItem, error)
	Find(req *request.GCartReq) (*model.Cart, error)
	CreateOrder(req *request.COrderReq) (*model.Order, error)
	UpdateOrderStatus(req *request.UOrderStatusReq) (*model.Order, error)
	DeleteCartItem(req *request.GCartItemReq) (bool, error)
}

type cartService struct {
	app          *app.Config
	cartRepo     repository.CartRepository
	cartItemRepo repository.CartItemRepository
	orderRepo    repository.OrderRepository
	productRepo  repository.ProductRepository
}

func NewCartService(appConfig *app.Config, cartRepo repository.CartRepository) CartService {
	cartItemRepo := repository.NewCartItemRepository(appConfig)
	orderRepo := repository.NewOrderRepository(appConfig)
	productRepo := repository.NewProductRepository(appConfig)
	return &cartService{
		appConfig,
		cartRepo,
		cartItemRepo,
		orderRepo,
		productRepo,
	}
}

func (s *cartService) AddToCart(item model.CartItem) (*model.CartItem, error) {
	return s.cartItemRepo.AddToCart(item)
}

func (s *cartService) Find(req *request.GCartReq) (*model.Cart, error) {
	return s.cartRepo.Find(req.ID)
}

func (s *cartService) CreateOrder(req *request.COrderReq) (*model.Order, error) {
	cart, err := s.cartRepo.Find(req.CartID)
	if err != nil {
		return nil, err
	}

	var total float64
	for _, item := range cart.Products {
		product, err := s.productRepo.Find(item.ProductID.String())
		if err != nil {
			return nil, err
		}
		total += product.Price * float64(item.Quantity)
	}

	m := model.Order{
		UserID: cart.UserID,
		Total:  total,
		Status: "pending",
	}

	order, err := s.orderRepo.Create(m)
	if err != nil {
		return nil, err
	}

	//ok, err := s.cartItemRepo.Deletes(cart.Products)
	//if err != nil || !ok {
	//	return nil, fmt.Errorf("error: can't deleted")
	//}

	return order, nil
}

func (s *cartService) UpdateOrderStatus(req *request.UOrderStatusReq) (*model.Order, error) {
	return s.orderRepo.UpdateOrderStatus(req)
}

func (s *cartService) DeleteCartItem(req *request.GCartItemReq) (bool, error) {
	return s.cartItemRepo.Delete(req.ID)
}
