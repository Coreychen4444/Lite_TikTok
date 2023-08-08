package repository

import (
	"github.com/Coreychen4444/Lite_TikTok/model"
	"gorm.io/gorm"
)

type DbRepository struct {
	db *gorm.DB
}

func NewDbRepository(db *gorm.DB) *DbRepository {
	return &DbRepository{db: db}
}

// 创建用户
func (r *DbRepository) CreateUsers(user *model.User) (*model.User, error) {
	result := r.db.Create(user)
	if result.Error != nil {
		return nil, result.Error
	}
	return user, nil
}

// 根据用户名获取用户
func (r *DbRepository) GetUserByName(username string) (*model.User, error) {
	var user model.User
	err := r.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, err
		}
		return nil, err
	}
	return &user, nil
}

// 根据用户id获取用户
func (r *DbRepository) GetUserById(id int64) (*model.User, error) {
	var user model.User
	err := r.db.Where("id = ?", id).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, err
		}
		return nil, err
	}
	return &user, nil
}
