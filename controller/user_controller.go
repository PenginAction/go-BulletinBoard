package controller

import (
	"net/http"

	"github.com/PenginAction/go-BulletinBoard/dto"
	"github.com/PenginAction/go-BulletinBoard/usecase"
	"github.com/labstack/echo/v4"
)

type IUserController interface {
	Signup(ctx echo.Context) error
	Login(ctx echo.Context) error
}

type userController struct {
	userUsecase usecase.IUserUsecase
}

func NewUserController(userUsecase usecase.IUserUsecase) IUserController {
	return &userController{userUsecase}
}

func (uc *userController) Signup(ctx echo.Context) error {
	var req dto.CreateUserRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}

	if err := ctx.Validate(req); err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}

	c := ctx.Request().Context()
	userRes, err := uc.userUsecase.SignUp(c, req)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusCreated, userRes)
}

func (uc *userController) Login(ctx echo.Context) error {
	var req dto.LoginRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}
	c := ctx.Request().Context()
	if err := ctx.Validate(req); err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}
	t, err := uc.userUsecase.Login(c, req)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, echo.Map{
		"token": t,
	})
}
