package service

import (
	"context"
	"task-manager/internal/repository"
	"time"
)

type AdminUserService struct {
	repository *repository.AdminUserRepository
}

func NewAdminUserService(repository *repository.AdminUserRepository) *AdminUserService {
	// 返回新结构体的地址；左侧 repository 是字段，右侧是函数参数。
	return &AdminUserService{
		repository: repository,
	}
}

type AdminUserReq struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	// CreatedAt 是创建时间；time.Time 是 time 包提供的时间类型。
	CreatedAt time.Time `json:"created_at"`
	// UpdatedAt 是最后更新时间；JSON 字段名为 updated_at。
	UpdatedAt time.Time `json:"updated_at"`
}

func (s *AdminUserService) Register(ctx context.Context, username string, possword string) (AdminUserReq, error) {

	adminuser, err := s.repository.Register(ctx, username, possword)
	if err != nil {
		return AdminUserReq{}, err
	}
	var reqAdminuser AdminUserReq
	reqAdminuser.ID = adminuser.ID
	reqAdminuser.Username = adminuser.Username
	reqAdminuser.CreatedAt = adminuser.CreatedAt
	reqAdminuser.UpdatedAt = adminuser.UpdatedAt
	return reqAdminuser, nil

}
