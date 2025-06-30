package repositories

import (
	"github/feroddev/challengeV3/internal/core"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	db.AutoMigrate(&core.User{})
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) Save(user *core.User) error {
	return r.db.Create(user).Error
}

func (r *UserRepository) FindByUsername(username string) (*core.User, error) {
	var user core.User
	result := r.db.Where("username = ?", username).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (r *UserRepository) FindByID(id uint) (*core.User, error) {
	var user core.User
	result := r.db.First(&user, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (r *UserRepository) Update(user *core.User) error {
	return r.db.Save(user).Error
}

func (r *UserRepository) Delete(id uint) error {
	return r.db.Delete(&core.User{}, id).Error
}

func (r *UserRepository) List() ([]core.User, error) {
	var users []core.User
	result := r.db.Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}
	return users, nil
}
