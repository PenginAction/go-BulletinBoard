package db

import (
	"context"
	"testing"
	"time"

	"github.com/PenginAction/go-BulletinBoard/utils"
	"github.com/stretchr/testify/require"
)

func CreateRandomPost(t *testing.T, user User) Post {
	arg := CreatePostParams{
		UserID: user.ID,
		Text:   utils.RandomString(9),
	}

	post, err := testQueries.CreatePost(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, post)

	require.Equal(t, arg.UserID, post.UserID)
	require.Equal(t, arg.Text, post.Text)

	require.NotZero(t, post.ID)
	require.NotZero(t, post.CreatedAt)

	return post
}

func TestCreatePost(t *testing.T) {
	user := createRandomUser(t)
	CreateRandomPost(t, user)
}

func TestGetPost(t *testing.T) {
	user := createRandomUser(t)
	post1 := CreateRandomPost(t, user)
	post2, err := testQueries.GetPost(context.Background(), post1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, post2)

	require.Equal(t, post1.UserID, post2.UserID)
	require.Equal(t, post1.Text, post2.Text)

	require.WithinDuration(t, post1.CreatedAt, post2.CreatedAt, time.Second)
}

func TestListPosts(t *testing.T) {
	user := createRandomUser(t)
	for i := 0; i < 10; i++ {
		CreateRandomPost(t, user)
	}

	arg := ListPostsParams{
		Limit:  5,
		Offset: 5,
	}

	posts, err := testQueries.ListPosts(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, posts, 5)

	for _, post := range posts {
		require.NotEmpty(t, post)
	}
}

func TestUpdatePost(t *testing.T) {
	user := createRandomUser(t)
	post1 := CreateRandomPost(t, user)

	arg := UpdatePostParams{
		ID:   post1.ID,
		Text: utils.RandomString(9),
	}

	post2, err := testQueries.UpdatePost(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, post2)

	require.Equal(t, post1.ID, post2.ID)
	require.Equal(t, post1.UserID, post2.UserID)
	require.Equal(t, arg.Text, post2.Text)
	require.WithinDuration(t, post1.CreatedAt, post2.CreatedAt, time.Second)
}

func TestDeletePost(t *testing.T) {
	user := createRandomUser(t)
	post1 := CreateRandomPost(t, user)
	err := testQueries.DeletePost(context.Background(), post1.ID)
	require.NoError(t, err)

	post2, err := testQueries.GetPost(context.Background(), post1.ID)
	require.Error(t, err)
	require.EqualError(t, err, err.Error())
	require.Empty(t, post2)
}
