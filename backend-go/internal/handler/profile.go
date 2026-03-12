package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/xyu/mdxy/internal/pkg/response"
	"github.com/xyu/mdxy/internal/service"
)

type ProfileHandler struct {
	service *service.ProfileService
}

func NewProfileHandler(service *service.ProfileService) *ProfileHandler {
	return &ProfileHandler{service: service}
}

// Get 获取个人信息
func (h *ProfileHandler) Get(c *gin.Context) {
	profile, err := h.service.GetProfile()
	if err != nil {
		response.InternalServerError(c, "获取个人信息失败")
		return
	}

	response.Success(c, profile)
}

// Update 更新个人信息（需要认证）
func (h *ProfileHandler) Update(c *gin.Context) {
	var req struct {
		Name     string `json:"name"`
		Title    string `json:"title"`
		Bio      string `json:"bio"`
		Avatar   string `json:"avatar"`
		Email    string `json:"email"`
		GitHub   string `json:"github"`
		LinkedIn string `json:"linkedin"`
		Twitter  string `json:"twitter"`
		Website  string `json:"website"`
		Skills   string `json:"skills"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误")
		return
	}

	profile, err := h.service.GetProfile()
	if err != nil {
		response.InternalServerError(c, "获取个人信息失败")
		return
	}

	// 更新字段
	profile.Name = req.Name
	profile.Title = req.Title
	profile.Bio = req.Bio
	profile.Avatar = req.Avatar
	profile.Email = req.Email
	profile.GitHub = req.GitHub
	profile.LinkedIn = req.LinkedIn
	profile.Twitter = req.Twitter
	profile.Website = req.Website
	profile.Skills = req.Skills

	if err := h.service.UpdateProfile(profile); err != nil {
		response.InternalServerError(c, "更新个人信息失败")
		return
	}

	response.SuccessWithMessage(c, "更新成功", profile)
}
