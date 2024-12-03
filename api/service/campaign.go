package service

import (
	"trinity/api/repository"
	"trinity/internal/model"
	"trinity/internal/request"
)

type CampaignService interface {
	Create(campaign model.Campaign) (*model.Campaign, error)
	Find(r *request.GCampaignReq) (*model.Campaign, error)
	List(r *request.LCampaignReq) ([]*model.Campaign, error)
	Count(r *request.LCampaignReq) (int64, error)
	Update(campaign model.Campaign) (*model.Campaign, error)
	Delete(r *request.GCampaignReq) (bool, error)
}

type campaignService struct {
	campaignRepo repository.CampaignRepository
}

func NewCampaignService(campaignRepo repository.CampaignRepository) CampaignService {
	return &campaignService{campaignRepo: campaignRepo}
}

func (s *campaignService) Create(campaign model.Campaign) (*model.Campaign, error) {
	return s.campaignRepo.Create(campaign)
}

func (s *campaignService) Find(r *request.GCampaignReq) (*model.Campaign, error) {
	return s.campaignRepo.Find(r.ID)
}

func (s *campaignService) List(r *request.LCampaignReq) ([]*model.Campaign, error) {
	return s.campaignRepo.List(r)
}

func (s *campaignService) Count(r *request.LCampaignReq) (int64, error) {
	return s.campaignRepo.Count(r)
}

func (s *campaignService) Update(user model.Campaign) (*model.Campaign, error) {
	return s.campaignRepo.Update(user)
}

func (s *campaignService) Delete(r *request.GCampaignReq) (bool, error) {
	return s.campaignRepo.Delete(r.ID)
}
