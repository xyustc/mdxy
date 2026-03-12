package repository

import (
	"gorm.io/gorm"

	"github.com/xyu/mdxy/internal/model"
)

type AdminRepository struct {
	db *gorm.DB
}

func NewAdminRepository(db *gorm.DB) *AdminRepository {
	return &AdminRepository{db: db}
}

func (r *AdminRepository) FindByUsername(username string) (*model.Admin, error) {
	var admin model.Admin
	err := r.db.Where("username = ?", username).First(&admin).Error
	return &admin, err
}

func (r *AdminRepository) UpdateLastLogin(id uint) error {
	return r.db.Model(&model.Admin{}).Where("id = ?", id).Update("last_login_at", gorm.Expr("CURRENT_TIMESTAMP")).Error
}
