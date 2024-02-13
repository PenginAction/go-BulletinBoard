package usecase

import (
	"context"

	db "github.com/PenginAction/go-BulletinBoard/db/sqlc"
	"github.com/PenginAction/go-BulletinBoard/dto"
	"github.com/PenginAction/go-BulletinBoard/utils"
)

type IUserUsecase interface {
	SignUp(c context.Context, req dto.CreateUserRequest) (dto.CreateUserResponse, error)
	Login(c context.Context, req dto.LoginRequest) (string, error)
}

type userUsecase struct {
	userRepository db.Querier
}

func NewUserUsecase(userRepository db.Querier) IUserUsecase {
	return &userUsecase{userRepository}
}

func (uu *userUsecase) SignUp(c context.Context, req dto.CreateUserRequest) (dto.CreateUserResponse, error) {
	hashPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return dto.CreateUserResponse{}, err
	}

	newUser := db.CreateUserParams{
		UserStrID: req.UserStrID,
		Email:     req.Email,
		Password:  hashPassword,
	}

	user, err := uu.userRepository.CreateUser(c, newUser)
	if err != nil {
		return dto.CreateUserResponse{}, err
	}

	rep := dto.CreateUserResponse{
		ID:        user.ID,
		UserStrID: user.UserStrID,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}

	return rep, nil
}

func (uu *userUsecase) Login(c context.Context, req dto.LoginRequest) (string, error) {
	user, err := uu.userRepository.GetUserByEmail(c, req.Email)
	if err != nil {
		return "", err
	}

	err = utils.CheckPassword(req.Password, user.Password)
	if err != nil {
		return "", err
	}

	token, err := utils.CreateValidToken(user.ID)
	if err != nil {
		return "", err
	}
	return token, nil
}
