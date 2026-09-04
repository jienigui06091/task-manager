package model

import "context"

// type Task struct 定义名为 Task 的结构体；结构体把多个相关字段组合成一个类型。
type AdminUser struct {
	// ID 是任务的数据库主键；int64 是 64 位有符号整数。
	ID int64 `json:"id"`
	// 反引号中的 json 标签规定 JSON 输出时的字段名为 title。
	PasswordHash string `json:"password_hash"`
	Username     string `json:"username"`
	Base
}

func (AdminUser) TableName() string {
	return "admin_user"
}

func CreateAdminUser(ctx context.Context, user *AdminUser) error {
	return DB.WithContext(ctx).Create(user).Error
}

func GetAdminUserByUsername(ctx context.Context, username string) (AdminUser, error) {
	var user AdminUser
	if err := DB.WithContext(ctx).Where("username = ?", username).First(&user).Error; err != nil {
		return AdminUser{}, err
	}
	return user, nil
}

func GetAdminUsers(ctx context.Context, page, pageSize int) ([]AdminUser, int64, error) {
	var total int64
	if err := DB.WithContext(ctx).Model(&AdminUser{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var users []AdminUser
	if err := DB.WithContext(ctx).
		Order("id DESC").
		Limit(pageSize).
		Offset(Offset(page, pageSize)).
		Find(&users).Error; err != nil {
		return nil, 0, err
	}
	return users, total, nil
}
