package usecase

import (
	"context"

	db "github.com/PenginAction/go-BulletinBoard/db/sqlc"
	"github.com/PenginAction/go-BulletinBoard/dto"
)

type IPostUsecase interface {
	CreatePost(c context.Context, req dto.CreatePostRequest) (dto.PostResponse, error)
	GetPostById(c context.Context, id uint) (dto.PostResponse, error)
	GetAllPosts(c context.Context, req dto.AllPostsRequest) ([]dto.PostResponse, error)
	UpdatePost(c context.Context, req dto.UpdatePostRequest) (dto.PostResponse, error)
	DeletePost(c context.Context, id uint) error
}

type postUsecase struct {
	postRepository db.Querier
}

func NewPostUsecase(postRepository db.Querier) IPostUsecase {
	return &postUsecase{postRepository}
}

func (pu *postUsecase) CreatePost(c context.Context, req dto.CreatePostRequest) (dto.PostResponse, error) {
	newPost := db.CreatePostParams{
		UserID: req.UserID,
		Text:   req.Text,
	}
	post, err := pu.postRepository.CreatePost(c, newPost)
	if err != nil {
		return dto.PostResponse{}, err
	}

	rep := dto.PostResponse{
		ID:        post.ID,
		UserID:    post.UserID,
		Text:      post.Text,
		CreatedAt: post.CreatedAt,
	}

	return rep, nil
}

func (pu *postUsecase) GetPostById(c context.Context, id uint) (dto.PostResponse, error) {
	post, err := pu.postRepository.GetPost(c, id)
	if err != nil {
		return dto.PostResponse{}, err
	}
	resPost := dto.PostResponse{
		ID:        post.ID,
		UserID:    post.UserID,
		Text:      post.Text,
		CreatedAt: post.CreatedAt,
	}
	return resPost, nil
}

func (pu *postUsecase) GetAllPosts(c context.Context, req dto.AllPostsRequest) ([]dto.PostResponse, error) {
	arg := db.ListPostsParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}
	posts, err := pu.postRepository.ListPosts(c, arg)
	if err != nil {
		return []dto.PostResponse{}, err
	}
	resPosts := []dto.PostResponse{}
	for _, v := range posts {
		p := dto.PostResponse{
			ID:        v.ID,
			UserID:    v.UserID,
			Text:      v.Text,
			CreatedAt: v.CreatedAt,
		}
		resPosts = append(resPosts, p)
	}
	return resPosts, nil
}

func (pu *postUsecase) UpdatePost(c context.Context, req dto.UpdatePostRequest) (dto.PostResponse, error) {
	renewPost := db.UpdatePostParams{
		ID:   req.ID,
		Text: req.Text,
	}
	post, err := pu.postRepository.UpdatePost(c, renewPost)
	if err != nil {
		return dto.PostResponse{}, err
	}
	resPost := dto.PostResponse{
		ID:        post.ID,
		UserID:    post.UserID,
		Text:      post.Text,
		CreatedAt: post.CreatedAt,
	}
	return resPost, nil
}

func (pu *postUsecase) DeletePost(c context.Context, id uint) error {
	if err := pu.postRepository.DeletePost(c, id); err != nil {
		return err
	}
	return nil
}
