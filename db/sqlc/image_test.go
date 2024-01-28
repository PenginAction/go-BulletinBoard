package db

import (
	"context"
	"testing"
	"time"

	"github.com/PenginAction/go-BulletinBoard/utils"
	"github.com/stretchr/testify/require"
)

func CreateRandomImage(t *testing.T, post Post) Image {
	arg := CreateImageParams{
		PostID:    post.ID,
		ImagePath: utils.RandomURL(),
	}

	image, err := testQueries.CreateImage(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, image)

	require.Equal(t, arg.PostID, image.PostID)
	require.Equal(t, arg.ImagePath, image.ImagePath)

	require.NotZero(t, image.ID)
	require.NotZero(t, image.CreatedAt)

	return image
}

func TestCreateImage(t *testing.T) {
	user := createRandomUser(t)
	post := CreateRandomPost(t, user)
	CreateRandomImage(t, post)
}

func TestGetImage(t *testing.T) {
	user := createRandomUser(t)
	post := CreateRandomPost(t, user)
	image1 := CreateRandomImage(t, post)
	image2, err := testQueries.GetImage(context.Background(), image1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, image2)

	require.Equal(t, image1.PostID, image2.PostID)
	require.Equal(t, image1.ImagePath, image2.ImagePath)

	require.WithinDuration(t, image1.CreatedAt, image2.CreatedAt, time.Second)
}

func TestListImages(t *testing.T) {
	user := createRandomUser(t)
	post := CreateRandomPost(t, user)
	for i := 0; i < 10; i++ {
		CreateRandomImage(t, post)
	}

	arg := ListImagesParams{
		PostID: post.ID,
		Limit:  5,
		Offset: 5,
	}

	images, err := testQueries.ListImages(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, images, 5)

	for _, image := range images {
		require.NotEmpty(t, image)
		require.Equal(t, arg.PostID, image.PostID)
	}
}

func TestUpdateImage(t *testing.T) {
	user := createRandomUser(t)
	post := CreateRandomPost(t, user)
	image1 := CreateRandomImage(t, post)

	arg := UpdateImageParams{
		ID:        image1.ID,
		ImagePath: utils.RandomURL(),
	}

	image2, err := testQueries.UpdateImage(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, image2)

	require.Equal(t, image1.ID, image2.ID)
	require.Equal(t, image1.PostID, image2.PostID)
	require.Equal(t, arg.ImagePath, image2.ImagePath)
	require.WithinDuration(t, image1.CreatedAt, image2.CreatedAt, time.Second)
}

func TestDeleteImage(t *testing.T) {
	user := createRandomUser(t)
	post := CreateRandomPost(t, user)
	image1 := CreateRandomImage(t, post)
	err := testQueries.DeleteImage(context.Background(), image1.ID)
	require.NoError(t, err)

	image2, err := testQueries.GetImage(context.Background(), image1.ID)
	require.Error(t, err)
	require.EqualError(t, err, err.Error())
	require.Empty(t, image2)
}
