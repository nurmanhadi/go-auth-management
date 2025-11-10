package repository

import (
	"auth-management/internal/entity"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}
func (r *UserRepository) Create(user *entity.User) error {
	return r.db.Save(user).Error
}
func (r *UserRepository) CountByUsername(username string) (int64, error) {
	var total int64
	err := r.db.Model(&entity.User{}).Where("username = ?", username).Count(&total).Error
	if err != nil {
		return 0, err
	}
	return total, nil
}
func (r *UserRepository) FindByUsername(username string) (*entity.User, error) {
	user := new(entity.User)
	err := r.db.Where("username = ?", username).First(user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}
