package usecase

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	mockdb "github.com/PenginAction/go-BulletinBoard/db/mock"
	db "github.com/PenginAction/go-BulletinBoard/db/sqlc"
	"github.com/PenginAction/go-BulletinBoard/dto"
	"github.com/PenginAction/go-BulletinBoard/utils"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

type eqCreateUserParamsMatcher struct {
	arg      db.CreateUserParams
	password string
}

func (e eqCreateUserParamsMatcher) Matches(x interface{}) bool {
	arg, ok := x.(db.CreateUserParams)
	if !ok {
		return false
	}

	err := utils.CheckPassword(e.password, arg.Password)
	if err != nil {
		return false
	}

	e.arg.Password = arg.Password
	return reflect.DeepEqual(e.arg, arg)
}

func (e eqCreateUserParamsMatcher) String() string {
	return fmt.Sprintf("matches arg %v and password %v", e.arg, e.password)
}

func EqCreateUserParams(arg db.CreateUserParams, password string) gomock.Matcher {
	return eqCreateUserParamsMatcher{arg, password}
}

func TestSignUp(t *testing.T) {
	user, password := RandomUser(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	arg := db.CreateUserParams{
		UserStrID: user.UserStrID,
		Email:     user.Email,
		Password:  password,
	}

	store := mockdb.NewMockStore(ctrl)
	store.EXPECT().
		CreateUser(gomock.Any(), EqCreateUserParams(arg, password)).
		Times(1).
		Return(user, nil)

	req := dto.CreateUserRequest{
		UserStrID: user.UserStrID,
		Email:     user.Email,
		Password:  password,
	}

	uu := NewUserUsecase(store)
	res, err := uu.SignUp(context.Background(), req)
	require.NoError(t, err)

	require.Equal(t, user.UserStrID, res.UserStrID)
	require.Equal(t, user.Email, res.Email)
}

func TestLogin(t *testing.T) {
	user, password := RandomUser(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mockdb.NewMockStore(ctrl)
	store.EXPECT().
		GetUserByEmail(gomock.Any(), gomock.Eq(user.Email)).
		Times(1).
		Return(user, nil)

	req := dto.LoginRequest{
		Email:    user.Email,
		Password: password,
	}

	uu := NewUserUsecase(store)
	token, err := uu.Login(context.Background(), req)
	require.NoError(t, err)
	require.NotEmpty(t, token)
}

func RandomUser(t *testing.T) (user db.User, password string) {
	password = utils.RandomString(6)
	hashedPassword, err := utils.HashPassword(password)
	require.NoError(t, err)

	user = db.User{
		UserStrID: utils.RandomUserStrID(),
		Email:     utils.RandomEmail(),
		Password:  hashedPassword,
	}
	return
}
