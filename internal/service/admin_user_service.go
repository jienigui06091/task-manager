package service

import (
	"context"
	"time"

	"golang.org/x/crypto/bcrypt"

	"task-manager/internal/auth"
	"task-manager/internal/model"
)

type AdminUserReq struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func RegisterAdminUser(ctx context.Context, username, password string) (AdminUserReq, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return AdminUserReq{}, err
	}

	user := model.AdminUser{
		Username:     username,
		PasswordHash: string(hash),
	}
	if err := model.CreateAdminUser(ctx, &user); err != nil {
		return AdminUserReq{}, err
	}

	return AdminUserReq{
		ID:        user.ID,
		Username:  user.Username,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

func LoginAdminUser(ctx context.Context, username, password string) (string, error) {
	user, err := model.GetAdminUserByUsername(ctx, username)
	if err != nil {
		return "", err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", err
	}

	return auth.GenerateJwtToken(username, user.ID)
}

func GetAdminUsers(ctx context.Context, page, pageSize int) ([]AdminUserReq, int64, error) {
	users, total, err := model.GetAdminUsers(ctx, page, pageSize)
	if err != nil {
		return nil, 0, err
	}

	result := make([]AdminUserReq, 0, len(users))
	for _, user := range users {
		result = append(result, AdminUserReq{
			ID:        user.ID,
			Username:  user.Username,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		})
	}
	return result, total, nil
}
