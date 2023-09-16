package jwthelper

import (
	"book-nest/config"
	"book-nest/internal/models/user"
	"errors"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	jwt5 "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func GenerateToken(input *user.User) (string, error) {
	claims := jwt5.MapClaims{
		"id":      input.Id,
		"name":    input.Name,
		"email":   input.Email,
		"phone":   input.Phone,
		"address": input.Address,
		"role":    input.Role,
		"exp":     time.Now().Add(time.Hour * time.Duration(config.Cfg.JwtConf.ExpiredTimeInHour)),
	}

	parseToken := jwt5.NewWithClaims(jwt5.SigningMethodHS256, claims)
	return parseToken.SignedString([]byte(config.Cfg.JwtConf.SecretKey))

}

func VerifyToken(ctx *gin.Context) (any, error) {
	headerToken := ctx.Request.Header.Get("Authorization")
	bearer := strings.HasPrefix(headerToken, "Bearer")

	if !bearer {
		return nil, errors.New("header token is not Bearer")
	}

	stringToken := strings.Split(headerToken, " ")[1]

	token, _ := jwt5.Parse(stringToken, func(t *jwt5.Token) (any, error) {
		if _, ok := t.Method.(*jwt5.SigningMethodHMAC); !ok {
			return nil, errors.New("error signing method HMAC")
		}
		return []byte(config.Cfg.JwtConf.SecretKey), nil
	})

	claims, ok := token.Claims.(jwt5.MapClaims)
	if !ok && !token.Valid {
		return nil, errors.New("invalid token")
	}

	expClaim, exists := claims["exp"]
	if !exists {
		return nil, errors.New("expire claim is missing")
	}

	expStr, ok := expClaim.(string)
	if !ok {
		return nil, errors.New("expire claim is not a valid type")
	}

	expTime, err := time.Parse(time.RFC3339, expStr)
	if err != nil {
		return nil, errors.New("error parsing expiration time")
	}

	if time.Now().After(expTime) {
		return nil, errors.New("token is expired")
	}

	return token.Claims.(jwt5.MapClaims), nil
}

func ParseToken(c *gin.Context) (*user.User, error) {
	userData := c.Value("userData").(jwt5.MapClaims)
	if userData == nil {
		return nil, errors.New("there is no user data provided")
	}
	userId, _ := uuid.Parse(userData["id"].(string))
	resp := &user.User{
		Id:      userId,
		Name:    userData["name"].(string),
		Email:   userData["email"].(string),
		Phone:   userData["phone"].(string),
		Address: userData["address"].(string),
		Role:    userData["role"].(string),
	}

	return resp, nil
}
