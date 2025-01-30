package services

import (
	"context"
	"strings"
	"time"
	"user-service/config"
	"user-service/constants"
	"user-service/domain/dto"
	"user-service/domain/models"
	"user-service/repositories"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	errConstants "user-service/constants/error"
)

type UserService struct {
	repository repositories.IRepositoryRegistry
}

type IUserService interface {
	Login(context.Context, *dto.LoginRequest) (*dto.LoginResponse, error)
	Register(context.Context, *dto.RegisterRequest) (*dto.RegisterRespose, error)
	Update(context.Context, *dto.UpdateRequest, string) (*dto.UserResponse, error)
	GetUserLogin(context.Context) (*dto.UserResponse, error)
	GetUserByUUID(context.Context, string) (*dto.UserResponse, error)
}

type Claims struct {
	User *dto.UserResponse
	jwt.RegisteredClaims
}

func NewUserService(repository repositories.IRepositoryRegistry) IUserService {
	return &UserService{repository: repository}
}

func (u *UserService) GetUserLogin(ctx context.Context) (*dto.UserResponse, error) {
	var (
		userLogin = ctx.Value(constants.UserLogin).(*dto.UserResponse)
		data      dto.UserResponse
	)

	data = dto.UserResponse{
		UUID:  userLogin.UUID,
		Name:  userLogin.Name,
		Email: userLogin.Email,
		Phone: userLogin.Phone,
		Role:  userLogin.Role,
	}

	return &data, nil
}

func (u *UserService) GetUserByUUID(ctx context.Context, uuid string) (*dto.UserResponse, error) {
	user, err := u.repository.GetUser().FindByUUID(ctx, uuid)
	if err != nil {
		return nil, err
	}

	data := dto.UserResponse{
		UUID:  user.UUID,
		Name:  user.Name,
		Email: user.Email,
		Phone: user.Phone,
	}

	return &data, nil
}

func (u *UserService) Login(ctx context.Context, req *dto.LoginRequest) (*dto.LoginResponse, error) {
	user, err := u.repository.GetUser().FindByEmail(ctx, req.Username)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return nil, err
	}

	expirationTime := time.Now().Add(time.Duration(config.Config.JwtExpirationTime) * time.Minute).Unix()
	data := &dto.UserResponse{
		UUID:  user.UUID,
		Name:  user.Name,
		Email: user.Email,
		Phone: user.Phone,
		Role:  strings.ToLower(user.Role.Code),
	}

	Claims := &Claims{
		User: data,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Unix(expirationTime, 0)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims)
	tokenString, err := token.SignedString([]byte(config.Config.JwtSecretKey))
	if err != nil {
		return nil, err
	}

	response := &dto.LoginResponse{
		User:  *data,
		Token: tokenString,
	}

	return response, nil
}

func (u *UserService) Register(ctx context.Context, req *dto.RegisterRequest) (*dto.RegisterRespose, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	if u.isEmailExist(ctx, req.Email) {
		return nil, errConstants.ErrEmailExists
	}

	if u.isPhoneExist(ctx, req.Phone) {
		return nil, errConstants.ErrPhoneExists
	}

	if req.Password != req.ConfirmPassword {
		return nil, errConstants.ErrPasswordDoesMatch
	}

	user, err := u.repository.GetUser().Register(ctx, &dto.RegisterRequest{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashedPassword),
		RoleID:   constants.Customer,
	})

	if err != nil {
		return nil, err
	}

	response := &dto.RegisterRespose{
		User: dto.UserResponse{
			UUID:  user.UUID,
			Name:  user.Name,
			Email: user.Email,
			Phone: user.Phone,
		},
	}

	return response, nil
}

func (u *UserService) isEmailExist(ctx context.Context, username string) bool {
	user, err := u.repository.GetUser().FindByEmail(ctx, username)
	if err != nil {
		return false
	}

	if user != nil {
		return true
	}

	return false
}

func (u *UserService) isPhoneExist(ctx context.Context, username string) bool {
	user, err := u.repository.GetUser().FindByPhone(ctx, username)
	if err != nil {
		return false
	}

	if user != nil {
		return true
	}

	return false
}

func (u *UserService) Update(ctx context.Context, req *dto.UpdateRequest, uuid string) (*dto.UserResponse, error) {
	var (
		password               string
		checkEmail, checkPhone *models.User
		hashedPassword         []byte
		user, userResult       *models.User
		err                    error
		data                   dto.UserResponse
	)

	user, err = u.repository.GetUser().FindByUUID(ctx, uuid)
	if err != nil {
		return nil, err
	}

	isEmailExist := u.isEmailExist(ctx, user.Email)
	if isEmailExist && user.Email != req.Email {
		checkEmail, err = u.repository.GetUser().FindByEmail(ctx, req.Email)
		if err != nil {
			return nil, err
		}
		if checkEmail != nil {
			return nil, errConstants.ErrEmailExists
		}
	}

	isPhoneExist := u.isPhoneExist(ctx, user.Phone)
	if isPhoneExist && user.Phone != req.Phone {
		checkPhone, err = u.repository.GetUser().FindByPhone(ctx, req.Phone)
		if err != nil {
			return nil, err
		}
		if checkPhone != nil {
			return nil, errConstants.ErrPhoneExists
		}
	}

	if req.Password != nil {
		if *req.Password != *req.ConfirmPassword {
			return nil, errConstants.ErrPasswordDoesMatch
		}
		hashedPassword, err = bcrypt.GenerateFromPassword([]byte(*req.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}

		password = string(hashedPassword)
	}

	userResult, err = u.repository.GetUser().Update(ctx, &dto.UpdateRequest{
		Name:     req.Name,
		Email:    req.Email,
		Password: &password,
		Phone:    req.Phone,
	}, uuid)

	if err != nil {
		return nil, err
	}

	data = dto.UserResponse{
		UUID:  userResult.UUID,
		Name:  userResult.Name,
		Email: userResult.Email,
		Phone: userResult.Phone,
	}

	return &data, nil
}
