package controller

import (
	"net/http"
	"strconv"

	"github.com/PenginAction/go-BulletinBoard/dto"
	"github.com/PenginAction/go-BulletinBoard/usecase"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type IPostController interface {
	CreatePost(ctx echo.Context) error
	GetPostById(ctx echo.Context) error
	GetAllPosts(ctx echo.Context) error
	UpdatePost(ctx echo.Context) error
	DeletePost(ctx echo.Context) error
}

type postController struct {
	postUsecase usecase.IPostUsecase
}

func NewPostController(pu usecase.IPostUsecase) IPostController {
	return &postController{pu}
}

func (pc *postController) CreatePost(ctx echo.Context) error {
	userValue := ctx.Get("user")
	if userValue == nil {
		return ctx.JSON(http.StatusUnauthorized, "Unauthorized")
	}
	user := userValue.(*jwt.Token)
	claims := user.Claims.(*dto.JwtCustomClaims)
	userId := claims.ID

	var req dto.CreatePostRequest
	req.UserID = userId
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}

	if err := ctx.Validate(req); err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}

	c := ctx.Request().Context()
	postRes, err := pc.postUsecase.CreatePost(c, req)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusCreated, postRes)
}

func (pc *postController) GetPostById(ctx echo.Context) error {
	userValue := ctx.Get("user")
	if userValue == nil {
		return ctx.JSON(http.StatusUnauthorized, "Unauthorized")
	}
	user := userValue.(*jwt.Token)
	claims := user.Claims.(*dto.JwtCustomClaims)
	userId := claims.ID
	Id := ctx.Param("postId")
	postId, _ := strconv.Atoi(Id)

	c := ctx.Request().Context()
	postRes, err := pc.postUsecase.GetPostById(c, uint(postId))
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}

	if postRes.UserID != userId {
		return ctx.JSON(http.StatusUnauthorized, "Unauthorized")
	}

	return ctx.JSON(http.StatusOK, postRes)
}

func (pc *postController) GetAllPosts(ctx echo.Context) error {
	var req dto.AllPostsRequest

	pageID, err := strconv.Atoi(ctx.QueryParam("page_id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}
	req.PageID = int32(pageID)

	pageSize, err := strconv.Atoi(ctx.QueryParam("page_size"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}
	req.PageSize = int32(pageSize)

	c := ctx.Request().Context()
	postRes, err := pc.postUsecase.GetAllPosts(c, req)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, postRes)
}

func (pc *postController) UpdatePost(ctx echo.Context) error {
	userValue := ctx.Get("user")
	if userValue == nil {
		return ctx.JSON(http.StatusUnauthorized, "Unauthorized")
	}
	user := userValue.(*jwt.Token)
	claims := user.Claims.(*dto.JwtCustomClaims)
	userId := claims.ID
	id := ctx.Param("postId")
	postId, _ := strconv.Atoi(id)

	var req dto.UpdatePostRequest
	req.ID = uint(postId)
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}

	c := ctx.Request().Context()
	postRes, err := pc.postUsecase.GetPostById(c, uint(postId))
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}

	if postRes.UserID != userId {
		return ctx.JSON(http.StatusUnauthorized, "Unauthorized")
	}

	newPost, err := pc.postUsecase.UpdatePost(c, req)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, newPost)
}

func (pc *postController) DeletePost(ctx echo.Context) error {
	userValue := ctx.Get("user")
	if userValue == nil {
		return ctx.JSON(http.StatusUnauthorized, "Unauthorized")
	}
	user := userValue.(*jwt.Token)
	claims := user.Claims.(*dto.JwtCustomClaims)
	userId := claims.ID
	id := ctx.Param("postId")
	postId, _ := strconv.Atoi(id)

	c := ctx.Request().Context()
	postRes, err := pc.postUsecase.GetPostById(c, uint(postId))
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}

	if postRes.UserID != userId {
		return ctx.JSON(http.StatusUnauthorized, "Unauthorized")
	}

	if err := pc.postUsecase.DeletePost(c, uint(postId)); err != nil {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}

	return ctx.NoContent(http.StatusNoContent)
}
