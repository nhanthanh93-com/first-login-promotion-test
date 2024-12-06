package service

import (
	"trinity/api/repository"
	"trinity/internal/model"
	"trinity/internal/request"
)

type CampaignService interface {
	Create(campaign model.Campaign) (*model.Campaign, error)
	Find(req *request.GCampaignReq) (*model.Campaign, error)
	List(req *request.LCampaignReq) ([]*model.Campaign, error)
	Count(req *request.LCampaignReq) (int64, error)
	Update(campaign model.Campaign) (*model.Campaign, error)
	Delete(req *request.GCampaignReq) (bool, error)
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

func (s *campaignService) Find(req *request.GCampaignReq) (*model.Campaign, error) {
	return s.campaignRepo.Find(req.ID)
}

func (s *campaignService) List(req *request.LCampaignReq) ([]*model.Campaign, error) {
	return s.campaignRepo.List(req)
}

func (s *campaignService) Count(req *request.LCampaignReq) (int64, error) {
	return s.campaignRepo.Count(req)
}

func (s *campaignService) Update(user model.Campaign) (*model.Campaign, error) {
	return s.campaignRepo.Update(user)
}

func (s *campaignService) Delete(req *request.GCampaignReq) (bool, error) {
	return s.campaignRepo.Delete(req.ID)
}
