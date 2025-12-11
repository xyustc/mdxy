package services

import (
	"mdxy-backend/config"
	"mdxy-backend/utils"
)

// AuthResult 认证结果
type AuthResult struct {
	Success bool   `json:"success"`
	Token   string `json:"token"`
	Message string `json:"message"`
}

// AuthenticateAdmin 验证管理员密码并返回token
func AuthenticateAdmin(password string) AuthResult {
	if password == config.AdminPassword {
		token, err := utils.GenerateJWT("admin")
		if err != nil {
			return AuthResult{
				Success: false,
				Token:   "",
				Message: "生成token失败",
			}
		}

		return AuthResult{
			Success: true,
			Token:   token,
			Message: "登录成功",
		}
	} else {
		return AuthResult{
			Success: false,
			Token:   "",
			Message: "密码错误",
		}
	}
}
