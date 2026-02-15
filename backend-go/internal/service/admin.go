package service

import (
	"errors"

	"golang.org/x/crypto/bcrypt"

	"github.com/xyu/mdxy/internal/pkg/jwt"
	"github.com/xyu/mdxy/internal/repository"
)

type AdminService struct {
	repo *repository.AdminRepository
}

func NewAdminService(repo *repository.AdminRepository) *AdminService {
	return &AdminService{repo: repo}
}

// Login 管理员登录
func (s *AdminService) Login(username, password string) (string, error) {
	admin, err := s.repo.FindByUsername(username)
	if err != nil {
		return "", errors.New("用户名或密码错误")
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(admin.PasswordHash), []byte(password)); err != nil {
		return "", errors.New("用户名或密码错误")
	}

	// 更新最后登录时间
	s.repo.UpdateLastLogin(admin.ID)

	// 生成 token
	token, err := jwt.GenerateToken(admin.ID, admin.Username)
	if err != nil {
		return "", errors.New("生成令牌失败")
	}

	return token, nil
}
