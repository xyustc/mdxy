package repository

import (
	"gorm.io/gorm"

	"github.com/xyu/mdxy/internal/model"
)

type ToolRepository struct {
	db *gorm.DB
}

func NewToolRepository(db *gorm.DB) *ToolRepository {
	return &ToolRepository{db: db}
}

func (r *ToolRepository) List(category string, visibleOnly bool) ([]model.Tool, error) {
	var tools []model.Tool
	query := r.db.Order("sort_order ASC, created_at DESC")
	if category != "" {
		query = query.Where("category = ?", category)
	}
	if visibleOnly {
		query = query.Where("is_visible = ?", true)
	}
	err := query.Find(&tools).Error
	return tools, err
}

func (r *ToolRepository) GetByID(id uint) (*model.Tool, error) {
	var tool model.Tool
	err := r.db.First(&tool, id).Error
	return &tool, err
}

func (r *ToolRepository) Create(tool *model.Tool) error {
	return r.db.Create(tool).Error
}

func (r *ToolRepository) Update(tool *model.Tool) error {
	return r.db.Save(tool).Error
}

func (r *ToolRepository) Delete(id uint) error {
	return r.db.Delete(&model.Tool{}, id).Error
}

func (r *ToolRepository) ListFeatured() ([]model.Tool, error) {
	var tools []model.Tool
	err := r.db.Where("is_featured = ? AND is_visible = ?", true, true).
		Order("sort_order ASC, created_at DESC").
		Find(&tools).Error
	return tools, err
}

func (r *ToolRepository) GetCategories() ([]string, error) {
	var categories []string
	err := r.db.Model(&model.Tool{}).Where("is_visible = ?", true).Distinct().Pluck("category", &categories).Error
	return categories, err
}
