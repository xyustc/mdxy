package repository

import (
	"gorm.io/gorm"

	"github.com/xyu/mdxy/internal/model"
)

type ProfileRepository struct {
	db *gorm.DB
}

func NewProfileRepository(db *gorm.DB) *ProfileRepository {
	return &ProfileRepository{db: db}
}

func (r *ProfileRepository) Get() (*model.Profile, error) {
	var profile model.Profile
	err := r.db.First(&profile).Error
	return &profile, err
}

func (r *ProfileRepository) Update(profile *model.Profile) error {
	return r.db.Save(profile).Error
}
