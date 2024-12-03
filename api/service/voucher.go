package service

import (
	"trinity/api/repository"
	"trinity/internal/model"
	"trinity/internal/request"
)

type VoucherService interface {
	Find(r *request.GVoucherReq) (*model.Voucher, error)
	List(r *request.LVoucherReq) ([]*model.Voucher, error)
	Count(r *request.LVoucherReq) (int64, error)
	Update(voucher model.Voucher) (*model.Voucher, error)
	Delete(r *request.GVoucherReq) (bool, error)
}

type voucherService struct {
	voucherRepo repository.VoucherRepository
}

func NewVoucherService(voucherRepo repository.VoucherRepository) VoucherService {
	return &voucherService{voucherRepo: voucherRepo}
}

func (s *voucherService) Find(r *request.GVoucherReq) (*model.Voucher, error) {
	return s.voucherRepo.Find(r.ID)
}

func (s *voucherService) List(r *request.LVoucherReq) ([]*model.Voucher, error) {
	return s.voucherRepo.List(r)
}

func (s *voucherService) Count(r *request.LVoucherReq) (int64, error) {
	return s.voucherRepo.Count(r)
}

func (s *voucherService) Update(user model.Voucher) (*model.Voucher, error) {
	return s.voucherRepo.Update(user)
}

func (s *voucherService) Delete(r *request.GVoucherReq) (bool, error) {
	return s.voucherRepo.Delete(r.ID)
}
