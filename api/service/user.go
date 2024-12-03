package service

import (
	"trinity/api/repository"
	"trinity/internal/model"
	"trinity/internal/request"
)

type UserService interface {
	Register(req *request.RUserReq) (*model.User, error)
	Create(user model.User) (*model.User, error)
	Find(r *request.GUserReq) (*model.User, error)
	List(r *request.LUserReq) ([]*model.User, error)
	Count(r *request.LUserReq) (int64, error)
	Update(user model.User) (*model.User, error)
	Delete(r *request.GUserReq) (bool, error)
}

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{userRepo: userRepo}
}

func (s *userService) Register(req *request.RUserReq) (*model.User, error) {
	return s.userRepo.Register(req)
}

func (s *userService) Create(m model.User) (*model.User, error) {
	return s.userRepo.Create(m)
}

func (s *userService) Find(r *request.GUserReq) (*model.User, error) {
	return s.userRepo.Find(r.ID)
}

func (s *userService) List(r *request.LUserReq) ([]*model.User, error) {
	return s.userRepo.List(r)
}

func (s *userService) Count(r *request.LUserReq) (int64, error) {
	return s.userRepo.Count(r)
}

func (s *userService) Update(user model.User) (*model.User, error) {
	return s.userRepo.Update(user)
}

func (s *userService) Delete(r *request.GUserReq) (bool, error) {
	return s.userRepo.Delete(r.ID)
}
