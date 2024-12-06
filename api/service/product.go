package service

import (
	"trinity/api/repository"
	"trinity/internal/model"
	"trinity/internal/request"
)

type ProductService interface {
	Create(product model.Product) (*model.Product, error)
	Find(req *request.GProductReq) (*model.Product, error)
	List(req *request.LProductReq) ([]*model.Product, error)
	Count(req *request.LProductReq) (int64, error)
	Update(product model.Product) (*model.Product, error)
	Delete(req *request.GProductReq) (bool, error)
}

type productService struct {
	productRepo repository.ProductRepository
}

func NewProductService(productRepo repository.ProductRepository) ProductService {
	return &productService{productRepo}
}

func (s *productService) Create(Product model.Product) (*model.Product, error) {
	return s.productRepo.Create(Product)
}

func (s *productService) Find(req *request.GProductReq) (*model.Product, error) {
	return s.productRepo.Find(req.ID)
}

func (s *productService) List(req *request.LProductReq) ([]*model.Product, error) {
	return s.productRepo.List(req)
}

func (s *productService) Count(req *request.LProductReq) (int64, error) {
	return s.productRepo.Count(req)
}

func (s *productService) Update(user model.Product) (*model.Product, error) {
	return s.productRepo.Update(user)
}

func (s *productService) Delete(req *request.GProductReq) (bool, error) {
	return s.productRepo.Delete(req.ID)
}
