package repository

import (
	"context"
	"errors"
	"user-service/domain/dto"
	"user-service/domain/models"

	"github.com/google/uuid"
	"gorm.io/gorm"

	commonErr "user-service/common/error"
	constantErr "user-service/constants/error"
)

type UserRepository struct {
	db *gorm.DB
}

type IUserRepository interface {
	Register(context.Context, *dto.RegisterRequest) (*models.User, error)
	Update(context.Context, *dto.UpdateRequest, string) (*models.User, error)
	FindByEmail(context.Context, string) (*models.User, error)
	FindByPhone(context.Context, string) (*models.User, error)
	FindByUUID(context.Context, string) (*models.User, error)
}

func NewUserRepository(db *gorm.DB) IUserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Register(ctx context.Context, req *dto.RegisterRequest) (*models.User, error) {
	user := models.User{
		UUID:     uuid.New(),
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
		Phone:    req.Phone,
		RoleID:   req.RoleID,
	}

	err := r.db.WithContext(ctx).Create(&user).Error
	if err != nil {
		return nil, commonErr.WrapError(constantErr.ErrSQLError)
	}

	return &user, nil
}

func (r *UserRepository) Update(ctx context.Context, req *dto.UpdateRequest, uuid string) (*models.User, error) {
	user := models.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: *req.Password,
		Phone:    req.Phone,
	}

	err := r.db.WithContext(ctx).Where("uuid = ?", uuid).Updates(&user).Error
	if err != nil {
		return nil, commonErr.WrapError(constantErr.ErrSQLError)
	}

	return &user, nil
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User

	err := r.db.WithContext(ctx).Preload("Role").Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, constantErr.ErrUserNotFound
		}
		return nil, commonErr.WrapError(constantErr.ErrSQLError)
	}

	return &user, nil
}

func (r *UserRepository) FindByPhone(ctx context.Context, phone string) (*models.User, error) {
	var user models.User

	err := r.db.WithContext(ctx).Preload("Role").Where("phone = ?", phone).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, constantErr.ErrUserNotFound
		}
		return nil, commonErr.WrapError(constantErr.ErrSQLError)
	}

	return &user, nil
}

func (r *UserRepository) FindByUUID(ctx context.Context, uuid string) (*models.User, error) {
	var user models.User

	err := r.db.WithContext(ctx).Preload("Role").Where("uuid = ?", uuid).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, constantErr.ErrUserNotFound
		}
		return nil, commonErr.WrapError(constantErr.ErrSQLError)
	}

	return &user, nil
}
