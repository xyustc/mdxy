package service

import (
	"errors"

	"github.com/xyu/mdxy/internal/model"
	"github.com/xyu/mdxy/internal/repository"
)

type ToolService struct {
	repo *repository.ToolRepository
}

func NewToolService(repo *repository.ToolRepository) *ToolService {
	return &ToolService{repo: repo}
}

func (s *ToolService) List(category string, visibleOnly bool) ([]model.Tool, error) {
	return s.repo.List(category, visibleOnly)
}

func (s *ToolService) GetByID(id uint) (*model.Tool, error) {
	return s.repo.GetByID(id)
}

func (s *ToolService) Create(tool *model.Tool) error {
	if tool.Name == "" {
		return errors.New("工具名称不能为空")
	}
	if tool.Type == "" {
		return errors.New("工具类型不能为空")
	}
	return s.repo.Create(tool)
}

func (s *ToolService) Update(tool *model.Tool) error {
	if tool.Name == "" {
		return errors.New("工具名称不能为空")
	}
	return s.repo.Update(tool)
}

func (s *ToolService) Delete(id uint) error {
	return s.repo.Delete(id)
}

func (s *ToolService) GetCategories() ([]string, error) {
	return s.repo.GetCategories()
}
