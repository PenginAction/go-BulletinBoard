package usecase

import (
	"context"
	"testing"

	mockdb "github.com/PenginAction/go-BulletinBoard/db/mock"
	db "github.com/PenginAction/go-BulletinBoard/db/sqlc"
	"github.com/PenginAction/go-BulletinBoard/dto"
	"github.com/PenginAction/go-BulletinBoard/utils"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestCreatePost(t *testing.T) {
	user, _ := RandomUser(t)
	post := RandomPost(user.ID)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	arg := db.CreatePostParams{
		UserID: post.UserID,
		Text:   post.Text,
	}

	store := mockdb.NewMockStore(ctrl)
	store.EXPECT().
		CreatePost(gomock.Any(), gomock.Eq(arg)).
		Times(1).
		Return(post, nil)

	req := dto.CreatePostRequest{
		UserID: post.UserID,
		Text:   post.Text,
	}

	pu := NewPostUsecase(store)
	res, err := pu.CreatePost(context.Background(), req)
	require.NoError(t, err)

	require.Equal(t, post.ID, res.ID)
	require.Equal(t, post.Text, req.Text)
}

func TestGetPost(t *testing.T) {
	user, _ := RandomUser(t)
	post := RandomPost(user.ID)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mockdb.NewMockStore(ctrl)
	store.EXPECT().
		GetPost(gomock.Any(), gomock.Eq(post.ID)).
		Times(1).
		Return(post, nil)

	pu := NewPostUsecase(store)
	res, err := pu.GetPostById(context.Background(), post.ID)
	require.NoError(t, err)

	require.Equal(t, post.ID, res.ID)
	require.Equal(t, post.Text, res.Text)
}

func TestGetAllPosts(t *testing.T) {
	user, _ := RandomUser(t)

	n := 5
	posts := make([]db.Post, n)
	for i := 0; i < n; i++ {
		posts[i] = RandomPost(user.ID)
	}

	arg := db.ListPostsParams{
		Limit:  int32(n),
		Offset: 0,
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mockdb.NewMockStore(ctrl)
	store.EXPECT().
		ListPosts(gomock.Any(), gomock.Eq(arg)).
		Times(1).
		Return(posts, nil)

	req := dto.AllPostsRequest{
		PageID:   1,
		PageSize: int32(n),
	}

	pu := NewPostUsecase(store)
	res, err := pu.GetAllPosts(context.Background(), req)
	require.NoError(t, err)

	require.Len(t, res, n)

	for i, post := range res {
		require.Equal(t, posts[i].ID, post.ID)
		require.Equal(t, posts[i].Text, post.Text)
	}
}

func TestUpdatePost(t *testing.T) {
	user, _ := RandomUser(t)
	post := RandomPost(user.ID)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	arg := db.UpdatePostParams{
		ID:   post.ID,
		Text: post.Text,
	}

	store := mockdb.NewMockStore(ctrl)
	store.EXPECT().
		UpdatePost(gomock.Any(), gomock.Eq(arg)).
		Times(1).
		Return(post, nil)

	req := dto.UpdatePostRequest{
		ID:   post.ID,
		Text: post.Text,
	}

	pu := NewPostUsecase(store)
	res, err := pu.UpdatePost(context.Background(), req)
	require.NoError(t, err)

	require.Equal(t, post.ID, res.ID)
	require.Equal(t, post.Text, req.Text)
}

func TestDeletePost(t *testing.T) {
	user, _ := RandomUser(t)
	post := RandomPost(user.ID)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mockdb.NewMockStore(ctrl)
	store.EXPECT().
		DeletePost(gomock.Any(), gomock.Eq(post.ID)).
		Times(1).
		Return(nil)

	pu := NewPostUsecase(store)
	err := pu.DeletePost(context.Background(), post.ID)
	require.NoError(t, err)
}

func RandomPost(userID uint) db.Post {
	post := db.Post{
		UserID: utils.RandomInt(1, 1000),
		Text:   utils.RandomString(15),
	}

	return post
}
