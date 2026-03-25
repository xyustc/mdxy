package repository

import (
	"time"

	"gorm.io/gorm"

	"github.com/xyu/mdxy/internal/model"
)

type CheatSheetRepository struct {
	db *gorm.DB
}

func NewCheatSheetRepository(db *gorm.DB) *CheatSheetRepository {
	return &CheatSheetRepository{db: db}
}

func (r *CheatSheetRepository) GetPublished(slug string) (*model.CheatSheetSnapshot, error) {
	var snapshot model.CheatSheetSnapshot
	err := r.db.Where("slug = ? AND status = ?", slug, "published").Order("published_at DESC, id DESC").First(&snapshot).Error
	return &snapshot, err
}

func (r *CheatSheetRepository) ListSnapshots(slug string, limit int) ([]model.CheatSheetSnapshot, error) {
	var snapshots []model.CheatSheetSnapshot
	query := r.db.Where("slug = ?", slug).Order("created_at DESC, id DESC")
	if limit > 0 {
		query = query.Limit(limit)
	}
	err := query.Find(&snapshots).Error
	return snapshots, err
}

func (r *CheatSheetRepository) GetByID(id uint) (*model.CheatSheetSnapshot, error) {
	var snapshot model.CheatSheetSnapshot
	err := r.db.First(&snapshot, id).Error
	return &snapshot, err
}

func (r *CheatSheetRepository) FindByHash(slug, contentHash string) (*model.CheatSheetSnapshot, error) {
	var snapshot model.CheatSheetSnapshot
	err := r.db.Where("slug = ? AND content_hash = ?", slug, contentHash).Order("id DESC").First(&snapshot).Error
	return &snapshot, err
}

func (r *CheatSheetRepository) Create(snapshot *model.CheatSheetSnapshot) error {
	return r.db.Create(snapshot).Error
}

func (r *CheatSheetRepository) Publish(slug string, id uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&model.CheatSheetSnapshot{}).
			Where("slug = ? AND status = ?", slug, "published").
			Updates(map[string]interface{}{"status": "draft", "published_at": nil}).Error; err != nil {
			return err
		}

		now := time.Now()
		return tx.Model(&model.CheatSheetSnapshot{}).
			Where("id = ? AND slug = ?", id, slug).
			Updates(map[string]interface{}{"status": "published", "published_at": &now, "sync_error": ""}).Error
	})
}
