package models

import (
	"fmt"
	"log/slog"

	"github.com/alexbsec/MiniMarketplace/src/db/config"
	"github.com/alexbsec/MiniMarketplace/src/db/models/utils"
	"github.com/alexbsec/MiniMarketplace/src/logging"
	"gorm.io/gorm"
)

type User struct {
	ID       uint    `gorm:"primaryKey"`
	Name     *string `gorm:"not null" json:"name"`
	Email    *string `gorm:"unique" json:"email"`
	Password *string `gorm:"not null" json:"password"`
    Role     *uint   `gorm:"not null" json:"role"`
}

type UserService struct {
	Service *config.Service
}

func (us *UserService) Create(user *User) error {
	if !us.isServiceRunning() {
		return fmt.Errorf("Cannot proceed because service is offline")
	}

	return models_utils.DoTransaction(us.Service, models_utils.CREATE, func(tx *gorm.DB) error {
		if err := tx.Create(user).Error; err != nil {
			return fmt.Errorf("failed to create user: %w", err)
		}

		return nil
	})
}

func (us *UserService) Fetch(id uint) (*User, error) {
	if !us.isServiceRunning() {
		return nil, fmt.Errorf("Cannot proceed because service is offline")
	}

	dbGorm, err := us.Service.Db()
	if err != nil {
		return nil, err
	}

	var user User
	res := dbGorm.First(&user, id)
	if res.Error != nil {
		if res.Error == gorm.ErrRecordNotFound {
			logging.Log.Error("Usuário não encontrado", slog.String("error", res.Error.Error()))
			return nil, res.Error
		}

		logging.Log.Error("Error while searching for user", slog.String("error", res.Error.Error()))
		return nil, res.Error
	}

	return &user, nil
}

func (us *UserService) Update(id uint, newUser *User) error {
	if !us.isServiceRunning() {
		return fmt.Errorf("Cannot proceed because service is offline")
	}

	return models_utils.DoTransaction(us.Service, models_utils.UPDATE, func(tx *gorm.DB) error {
		var user User
		if err := tx.First(&user, id).Error; err != nil {
			return fmt.Errorf("user with id %d not found: %w", id, err)
		}

		if err := tx.Model(&user).Updates(newUser).Error; err != nil {
			return fmt.Errorf("failed to update user with id %d: %w", id, err)
		}

		return nil
	})
}

func (us *UserService) Delete(id uint) error {
	if !us.isServiceRunning() {
		return fmt.Errorf("Cannot proceed because service is offline")
	}

	return models_utils.DoTransaction(us.Service, models_utils.DELETE, func(tx *gorm.DB) error {
		if err := tx.Delete(&User{}, id).Error; err != nil {
			return fmt.Errorf("failed to delete user: %w", err)
		}

		return nil
	})
}

func (us *UserService) FetchUserByEmail(email *string) (*User, error) {
 	if !us.isServiceRunning() {
		return nil, fmt.Errorf("Cannot proceed because service is offline")
	}

	dbGorm, err := us.Service.Db()
	if err != nil {
		return nil, err
	}

    var user User
    res := dbGorm.Where("email = ?", email).First(&user)
	if res.Error != nil {
		if res.Error == gorm.ErrRecordNotFound {
			logging.Log.Error("Usuário não encontrado", slog.String("error", res.Error.Error()))
			return nil, res.Error
		}

		logging.Log.Error("Error while searching for user", slog.String("error", res.Error.Error()))
		return nil, res.Error
	}

    return &user, nil
}

func (us *UserService) FetchPassword(id uint) (*string, error) {
	if !us.isServiceRunning() {
		return nil, fmt.Errorf("Cannot proceed because service is offline")
	}

	dbGorm, err := us.Service.Db()
	if err != nil {
		return nil, err
	}

	var user User 
	res := dbGorm.Select("password").First(&user, id)
	if res.Error != nil {
		if res.Error == gorm.ErrRecordNotFound {
			logging.Log.Error("Usuário não encontrado", slog.String("error", res.Error.Error()))
			return nil, res.Error
		}

		logging.Log.Error("Error while searching for user", slog.String("error", res.Error.Error()))
		return nil, res.Error
	}

    if user.Password == nil {
        return nil, fmt.Errorf("password is null for user with id %d", id)    
    } 

    return user.Password, nil
}

func (us *UserService) FetchEmail(id uint) (*string, error) {
	if !us.isServiceRunning() {
		return nil, fmt.Errorf("Cannot proceed because service is offline")
	}

	dbGorm, err := us.Service.Db()
	if err != nil {
		return nil, err
	}

	var user models_utils.MinifiedUser
	res := dbGorm.Select("email").First(&user, id)
	if res.Error != nil {
		if res.Error == gorm.ErrRecordNotFound {
			logging.Log.Error("Usuário não encontrado", slog.String("error", res.Error.Error()))
			return nil, res.Error
		}

		logging.Log.Error("Error while searching for user", slog.String("error", res.Error.Error()))
		return nil, res.Error
	}

    if user.Email == nil {
        return nil, fmt.Errorf("email is null for user with id %d", id) 
    }

    return user.Email, nil
}

func (us *UserService) CheckEmailExists(email *string) (bool, error) {
	if !us.isServiceRunning() {
		return false, fmt.Errorf("Cannot proceed because service is offline")
	}

	dbGorm, err := us.Service.Db()
	if err != nil {
		return false, err
	}

    var exists bool
    err = dbGorm.Raw("SELECT EXISTS(SELECT 1 FROM users WHERE email = ?)", email).
            Scan(&exists).Error

    if err != nil {
        return false, err
    }

    return exists, nil
}

func (us *UserService) isServiceRunning() bool {
	if us.Service == nil {
		logging.Log.Error("User Service is not initialized! Aborting")
	}

	return us.Service != nil
}
