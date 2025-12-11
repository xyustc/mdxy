package utils

import (
	"time"

	"mdxy-backend/config"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// Claims JWT声明结构
type Claims struct {
	Role string `json:"role"`
	jwt.RegisteredClaims
}

// HashPassword 哈希密码
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPasswordHash 验证密码哈希
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// GenerateJWT 生成JWT token
func GenerateJWT(role string) (string, error) {
	// 设置过期时间
	expirationTime := time.Now().Add(time.Hour * time.Duration(config.JWTExpireHours))

	// 创建声明
	claims := &Claims{
		Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	// 创建token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 签名token
	tokenString, err := token.SignedString([]byte(config.JWTSecretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateJWT 验证JWT token
func ValidateJWT(tokenString string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.JWTSecretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, jwt.ErrTokenUnverifiable
	}

	return claims, nil
}
