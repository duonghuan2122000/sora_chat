package service

import (
	"context"
	"encoding/json"
	"fmt"
	"sora_chat/internal/database"
	sora_errors "sora_chat/internal/errors"
	"sora_chat/internal/model"
	"sora_chat/internal/repository"
	"sora_chat/internal/util"
	"time"
)

// Dto req đăng nhập
type LoginUserReqDto struct {
	// Luồng cấp quyền
	GrantType model.GrantType `json:"grantType"`

	// Username
	Username string `json:"username"`

	// mật khẩu
	Password string `json:"password"`
}

// Dto res đăng nhập
type LoginUserResDto struct {
	// access token
	AccessToken string `json:"accessToken,omitempty"`

	// Loại access token
	TokenType string `json:"tokenType,omitempty"`

	// Thời gian hiệu lực của token, tính bằng giây
	ExpiresIn int `json:"expiresIn,omitempty"`
}

type TokenInfo struct {
	// UserId
	ID string `json:"sub"`
	// UserName
	Username string `json:"user_name"`
	// Tên
	FirstName string `json:"first_name"`
	// Họ và tên đệm
	LastName string `json:"last_name"`
	// Thời gian hết hạn
	ExpiresAt time.Time `json:"exp"`
}

type CreateUserReqDto struct {
	// Username
	Username string `json:"username"`

	// mật khẩu
	Password string `json:"password"`

	// Tên
	FirstName string `json:"firstName"`

	// Họ và tên đệm
	LastName string `json:"lastName"`
}

type CreateUserResDto struct {
	UserId string `json:"userId"`
}

type SessionUserInfoResDto struct {
	// UserId
	UserId string `json:"userId"`
	// UserName
	Username string `json:"userName"`
	// Tên
	FirstName string `json:"firstName"`
	// Họ và tên đệm
	LastName string `json:"lastName"`
	// Thời gian hết hạn
	ExpiresAt time.Time `json:"expiresAt"`
}

func NewLoginUserResDto() *LoginUserResDto {
	return &LoginUserResDto{
		TokenType: "Bearer",
		ExpiresIn: 86400,
	}
}

type UserService interface {

	// Tạo user
	CreateUser(payload CreateUserReqDto) (*CreateUserResDto, error)

	// Thực hiện cấp quyền truy cập
	GrantToken(payload LoginUserReqDto) (*LoginUserResDto, error)

	// lấy thông tin user của phiên hiện tại
	GetCurrentUserInfo(token string, tokenInfo *TokenInfo) (*SessionUserInfoResDto, error)
}

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{userRepo: userRepo}
}

func (s *userService) FindByEmail(email string) (*model.User, error) {
	user, err := s.userRepo.FindByEmail(email)
	return user, err
}

// Lưu thông tin token vào Redis
func saveTokenToRedis(key string, tokenInfo TokenInfo) error {
	dataJson, err := json.Marshal(tokenInfo)
	if err != nil {
		return err
	}
	duration := time.Until(tokenInfo.ExpiresAt)
	return database.RedisDB.Set(context.Background(), "sora:"+key, dataJson, duration).Err()
}

// Tạo user
func (s *userService) CreateUser(payload CreateUserReqDto) (*CreateUserResDto, error) {
	passwordHashed, _ := util.HashPassword(payload.Password)
	var user = model.User{
		Username:       payload.Username,
		Email:          payload.Username,
		FirstName:      payload.FirstName,
		LastName:       payload.LastName,
		PasswordHashed: passwordHashed,
	}
	userId, err := s.userRepo.Insert(user)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	result := &CreateUserResDto{
		UserId: *userId,
	}
	return result, nil
}

// Thực hiện cấp quyền truy cập
func (s *userService) GrantToken(payload LoginUserReqDto) (*LoginUserResDto, error) {
	var result *LoginUserResDto
	switch payload.GrantType {
	case model.GrantType_Password:
		user, err := s.userRepo.FindByEmail(payload.Username)
		if err != nil {
			return nil, &sora_errors.LogicError{
				Code:    "001",
				Message: "User/Pass không hợp lệ",
			}
		}
		validPassword := util.VerifyPassword(payload.Password, user.PasswordHashed)
		if !validPassword {
			return nil, &sora_errors.LogicError{
				Code:    "001",
				Message: "User/Pass không hợp lệ",
			}
		}
		tokenInfo := TokenInfo{
			ID:        user.ID,
			Username:  user.Username,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			ExpiresAt: time.Now().UTC().Add(1 * time.Hour),
		}
		result = NewLoginUserResDto()
		token, err := util.GenerateReferenceToken()
		if err != nil {
			return nil, err
		}
		err = saveTokenToRedis(token, tokenInfo)
		if err != nil {
			return nil, err
		}
		result.AccessToken = token
		return result, nil
	default:
		return nil, &sora_errors.NotSupportedError{}
	}
}

// lấy thông tin user của phiên hiện tại
func (s *userService) GetCurrentUserInfo(token string, tokenInfo *TokenInfo) (*SessionUserInfoResDto, error) {
	result := &SessionUserInfoResDto{
		UserId:    tokenInfo.ID,
		Username:  tokenInfo.Username,
		FirstName: tokenInfo.FirstName,
		LastName:  tokenInfo.LastName,
		ExpiresAt: tokenInfo.ExpiresAt,
	}
	return result, nil
}
