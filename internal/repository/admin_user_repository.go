package repository

import (
	"context"
	"task-manager/internal/model"

	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type AdminUserRepository struct {
	// db 是指向连接池的指针；* 表示保存地址而不是复制整个连接池。
	db *pgxpool.Pool
}

func NewAdminUseRepository(db *pgxpool.Pool) *AdminUserRepository {
	// &TaskRepository{...} 创建结构体并取其地址，返回 *TaskRepository 指针。
	return &AdminUserRepository{
		// 左侧 db 是字段名，右侧 db 是函数参数。
		db: db,
	}
}

func (r *AdminUserRepository) Register(cxt context.Context, username string, password string) (model.AdminUser, error) {
	newVar, newVar1 := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if newVar1 != nil {
		return model.AdminUser{}, newVar1
	}
	var admin model.AdminUser
	err := r.db.QueryRow(cxt, `
	insert into admin_user(username,password_hash)values(
	$1,$2
	)RETURNING id, username, password_hash, created_at, updated_at
	`, username, newVar).Scan(
		&admin.ID,
		&admin.Username,
		&admin.PasswordHash,
		&admin.CreatedAt,
		&admin.UpdatedAt,
	)
	if err != nil {
		return model.AdminUser{}, err

	}
	return admin, nil
}

func (r *AdminUserRepository) GetByUsername(ctx context.Context, username string) (model.AdminUser, error) {
	var user model.AdminUser
	err:=r.db.QueryRow(ctx, `
	select * from admin_user
	where username = $1
	`, username).Scan(
		&user.ID,
		&user.Username,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err!=nil {
		return model.AdminUser{},err
	}
	return user,nil
}
