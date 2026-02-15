package service

import (
	"github.com/xyu/mdxy/internal/model"
	"github.com/xyu/mdxy/internal/repository"
)

type ProfileService struct {
	repo *repository.ProfileRepository
}

func NewProfileService(repo *repository.ProfileRepository) *ProfileService {
	return &ProfileService{repo: repo}
}

func (s *ProfileService) GetProfile() (*model.Profile, error) {
	return s.repo.Get()
}

func (s *ProfileService) UpdateProfile(profile *model.Profile) error {
	return s.repo.Update(profile)
}
