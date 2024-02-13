package controller

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/PenginAction/go-BulletinBoard/dto"
	mock_usecase "github.com/PenginAction/go-BulletinBoard/usecase/mock"
	"github.com/PenginAction/go-BulletinBoard/utils"
	"github.com/go-playground/validator"
	"github.com/golang-jwt/jwt/v5"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
)

func TestCreatePost(t *testing.T) {
	cases := []struct {
		name          string
		requestBody   map[string]interface{}
		buildStubs    func(pu *mock_usecase.MockIPostUsecase, requestBody map[string]interface{})
		checkResponse func(recoder *httptest.ResponseRecorder)
	}{
		{
			name: "valid request",
			requestBody: map[string]interface{}{
				"user_id": utils.RandomInt(1, 100),
				"text":    utils.RandomString(6),
			},
			buildStubs: func(pu *mock_usecase.MockIPostUsecase, requestBody map[string]interface{}) {
				expectedPost := dto.CreatePostRequest{
					UserID: requestBody["user_id"].(uint),
					Text:   requestBody["text"].(string),
				}
				pu.EXPECT().
					CreatePost(context.Background(), expectedPost).
					Times(1).
					Return(dto.PostResponse{}, nil)
			},
			checkResponse: func(rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusCreated, rec.Code)
			},
		},
		{
			name:        "no user info",
			requestBody: map[string]interface{}{},
			buildStubs: func(pu *mock_usecase.MockIPostUsecase, requestBody map[string]interface{}) {
			},
			checkResponse: func(rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, rec.Code)
			},
		},
		{
			name: "request bind error",
			requestBody: map[string]interface{}{
				"user_id": "invalid user id",
				"text":    utils.RandomString(6),
			},
			buildStubs: func(pu *mock_usecase.MockIPostUsecase, requestBody map[string]interface{}) {
			},
			checkResponse: func(rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, rec.Code)
			},
		},
		{
			name: "validation error",
			requestBody: map[string]interface{}{
				"user_id": utils.RandomInt(1, 100),
				"text":    "",
			},
			buildStubs: func(pu *mock_usecase.MockIPostUsecase, requestBody map[string]interface{}) {
			},
			checkResponse: func(rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, rec.Code)
			},
		},
		{
			name: "internal server error",
			requestBody: map[string]interface{}{
				"user_id": utils.RandomInt(1, 100),
				"text":    utils.RandomString(6),
			},
			buildStubs: func(pu *mock_usecase.MockIPostUsecase, requestBody map[string]interface{}) {
				expectedPost := dto.CreatePostRequest{
					UserID: requestBody["user_id"].(uint),
					Text:   requestBody["text"].(string),
				}
				pu.EXPECT().
					CreatePost(context.Background(), expectedPost).
					Times(1).
					Return(dto.PostResponse{}, errors.New("internal server error"))
			},
			checkResponse: func(rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, rec.Code)
			},
		},
	}

	e := echo.New()
	e.Validator = &utils.CustomValidator{Validator: validator.New()}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	pu := mock_usecase.NewMockIPostUsecase(ctrl)
	pc := NewPostController(pu)

	for i := range cases {
		tc := cases[i]
		t.Run(tc.name, func(t *testing.T) {
			tc.buildStubs(pu, tc.requestBody)

			body, _ := json.Marshal(tc.requestBody)

			var userID uint
			if id, ok := tc.requestBody["user_id"].(uint); ok {
				userID = id
			}

			token, err := utils.CreateValidToken(userID)
			require.NoError(t, err)

			url := "/posts/"
			req := httptest.NewRequest(http.MethodPost, url, bytes.NewReader(body))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %v", token))
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			if tc.name != "no user info" {
				user := &jwt.Token{Claims: &dto.JwtCustomClaims{ID: userID}}
				c.Set("user", user)
			}
			err = pc.CreatePost(c)
			require.NoError(t, err)

			tc.checkResponse(rec)
		})
	}
}

func TestGetPostById(t *testing.T) {
	cases := []struct {
		name            string
		postID          uint
		requestBody     map[string]interface{}
		expectedPostRes dto.PostResponse
		buildStubs      func(pu *mock_usecase.MockIPostUsecase, postID uint, requestBody map[string]interface{}, expectedPostRes dto.PostResponse)
		checkResponse   func(recoder *httptest.ResponseRecorder)
	}{
		{
			name:   "valid request",
			postID: utils.RandomInt(1, 100),
			requestBody: map[string]interface{}{
				"user_id": utils.RandomInt(1, 100),
				"text":    utils.RandomString(6),
			},
			expectedPostRes: dto.PostResponse{
				ID:     utils.RandomInt(1, 100),
				UserID: utils.RandomInt(1, 100),
				Text:   utils.RandomString(6),
			},
			buildStubs: func(pu *mock_usecase.MockIPostUsecase, postID uint, requestBody map[string]interface{}, expectedPostRes dto.PostResponse) {
				pu.EXPECT().
					GetPostById(context.Background(), postID).
					Times(1).
					Return(expectedPostRes, nil)
			},
			checkResponse: func(rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, rec.Code)
			},
		},
		{
			name:   "no user info",
			postID: utils.RandomInt(1, 100),
			requestBody: map[string]interface{}{
				"user_id": utils.RandomInt(1, 100),
				"text":    utils.RandomString(6),
			},
			expectedPostRes: dto.PostResponse{
				ID:     utils.RandomInt(1, 100),
				UserID: utils.RandomInt(1, 100),
				Text:   utils.RandomString(6),
			},
			buildStubs: func(pu *mock_usecase.MockIPostUsecase, postID uint, requestBody map[string]interface{}, expectedPostRes dto.PostResponse) {
			},
			checkResponse: func(rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, rec.Code)
			},
		},
		{
			name:   "internal server error",
			postID: utils.RandomInt(1, 100),
			requestBody: map[string]interface{}{
				"user_id": utils.RandomInt(1, 100),
				"text":    utils.RandomString(6),
			},
			expectedPostRes: dto.PostResponse{
				ID:     utils.RandomInt(1, 100),
				UserID: utils.RandomInt(1, 100),
				Text:   utils.RandomString(6),
			},
			buildStubs: func(pu *mock_usecase.MockIPostUsecase, postID uint, requestBody map[string]interface{}, expectedPostRes dto.PostResponse) {
				pu.EXPECT().
					GetPostById(context.Background(), postID).
					Times(1).
					Return(expectedPostRes, errors.New("internal server error"))
			},
			checkResponse: func(rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, rec.Code)
			},
		},
		{
			name:   "post user id does not match token user id",
			postID: utils.RandomInt(1, 100),
			requestBody: map[string]interface{}{
				"user_id": utils.RandomInt(1, 100),
				"text":    utils.RandomString(6),
			},
			expectedPostRes: dto.PostResponse{
				ID:     utils.RandomInt(1, 100),
				UserID: utils.RandomInt(1, 100) + 1,
				Text:   utils.RandomString(6),
			},
			buildStubs: func(pu *mock_usecase.MockIPostUsecase, postID uint, requestBody map[string]interface{}, expectedPostRes dto.PostResponse) {
				pu.EXPECT().
					GetPostById(context.Background(), postID).
					Times(1).
					Return(expectedPostRes, nil)
			},
			checkResponse: func(rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, rec.Code)
			},
		},
	}
	e := echo.New()
	e.Validator = &utils.CustomValidator{Validator: validator.New()}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	pu := mock_usecase.NewMockIPostUsecase(ctrl)
	pc := NewPostController(pu)

	for i := range cases {
		tc := cases[i]
		t.Run(tc.name, func(t *testing.T) {
			tc.buildStubs(pu, tc.postID, tc.requestBody, tc.expectedPostRes)

			token, err := utils.CreateValidToken(tc.expectedPostRes.UserID)
			require.NoError(t, err)

			url := fmt.Sprintf("/posts/%d", tc.postID)
			req := httptest.NewRequest(http.MethodGet, url, nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %v", token))
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("postId")
			c.SetParamValues(strconv.Itoa(int(tc.postID)))

			if tc.name != "no user info" {
				user := &jwt.Token{Claims: &dto.JwtCustomClaims{ID: tc.expectedPostRes.UserID}}
				c.Set("user", user)
			}

			if tc.name == "post user id does not match token user id" {
				user := &jwt.Token{Claims: &dto.JwtCustomClaims{ID: tc.expectedPostRes.UserID + 1}}
				c.Set("user", user)
			}

			err = pc.GetPostById(c)
			require.NoError(t, err)

			tc.checkResponse(rec)
		})
	}
}

func TestGetAllPosts(t *testing.T) {
	cases := []struct {
		name          string
		requestBody   map[string]interface{}
		buildStubs    func(pu *mock_usecase.MockIPostUsecase, requestBody map[string]interface{})
		checkResponse func(recoder *httptest.ResponseRecorder)
	}{
		{
			name: "valid request",
			requestBody: map[string]interface{}{
				"user_id": utils.RandomInt(1, 100),
				"limit":   int32(utils.RandomInt(1, 100)),
				"offset":  int32(utils.RandomInt(1, 100)),
			},
			buildStubs: func(pu *mock_usecase.MockIPostUsecase, requestBody map[string]interface{}) {
				expectedReq := dto.AllPostsRequest{
					UserID: requestBody["user_id"].(uint),
					Limit:  requestBody["limit"].(int32),
					Offset: requestBody["offset"].(int32),
				}
				pu.EXPECT().
					GetAllPosts(context.Background(), expectedReq).
					Times(1).
					Return([]dto.PostResponse{}, nil)
			},
			checkResponse: func(rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, rec.Code)
			},
		},
		{
			name: "no user info",
			requestBody: map[string]interface{}{
				"user_id": utils.RandomInt(1, 100),
				"limit":   int32(utils.RandomInt(1, 100)),
				"offset":  int32(utils.RandomInt(1, 100)),
			},
			buildStubs: func(pu *mock_usecase.MockIPostUsecase, requestBody map[string]interface{}) {
			},
			checkResponse: func(rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, rec.Code)
			},
		},
		{
			name: "request bind error",
			requestBody: map[string]interface{}{
				"user_id": "invalid user id",
				"limit":   int32(utils.RandomInt(1, 100)),
				"offset":  int32(utils.RandomInt(1, 100)),
			},
			buildStubs: func(pu *mock_usecase.MockIPostUsecase, requestBody map[string]interface{}) {
			},
			checkResponse: func(rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, rec.Code)
			},
		},
		{
			name: "validation error",
			requestBody: map[string]interface{}{
				"user_id": utils.RandomInt(1, 100),
				"limit":   "",
				"offset":  "",
			},
			buildStubs: func(pu *mock_usecase.MockIPostUsecase, requestBody map[string]interface{}) {
			},
			checkResponse: func(rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, rec.Code)
			},
		},
		{
			name: "internal server error",
			requestBody: map[string]interface{}{
				"user_id": utils.RandomInt(1, 100),
				"limit":   int32(utils.RandomInt(1, 100)),
				"offset":  int32(utils.RandomInt(1, 100)),
			},
			buildStubs: func(pu *mock_usecase.MockIPostUsecase, requestBody map[string]interface{}) {
				expectedReq := dto.AllPostsRequest{
					UserID: requestBody["user_id"].(uint),
					Limit:  requestBody["limit"].(int32),
					Offset: requestBody["offset"].(int32),
				}
				pu.EXPECT().
					GetAllPosts(context.Background(), expectedReq).
					Times(1).
					Return([]dto.PostResponse{}, errors.New("internal server error"))
			},
			checkResponse: func(rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, rec.Code)
			},
		},
	}
	e := echo.New()
	e.Validator = &utils.CustomValidator{Validator: validator.New()}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	pu := mock_usecase.NewMockIPostUsecase(ctrl)
	pc := NewPostController(pu)

	for i := range cases {
		tc := cases[i]
		t.Run(tc.name, func(t *testing.T) {
			tc.buildStubs(pu, tc.requestBody)

			body, err := json.Marshal(tc.requestBody)
			require.NoError(t, err)

			var userID uint
			if id, ok := tc.requestBody["user_id"].(uint); ok {
				userID = id
			}
			token, err := utils.CreateValidToken(userID)
			require.NoError(t, err)

			url := "/posts/"
			req := httptest.NewRequest(http.MethodGet, url, bytes.NewReader(body))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %v", token))
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			if tc.name != "no user info" {
				user := &jwt.Token{Claims: &dto.JwtCustomClaims{ID: userID}}
				c.Set("user", user)
			}

			err = pc.GetAllPosts(c)
			require.NoError(t, err)

			tc.checkResponse(rec)
		})
	}
}

func TestUpdatePost(t *testing.T) {
	cases := []struct {
		name            string
		postID          uint
		requestBody     map[string]interface{}
		expectedPostRes dto.PostResponse
		buildStubs      func(pu *mock_usecase.MockIPostUsecase, postID uint, requestBody map[string]interface{}, expectedPostRes dto.PostResponse)
		checkResponse   func(recoder *httptest.ResponseRecorder)
	}{
		{
			name:   "valid request",
			postID: utils.RandomInt(1, 100),
			requestBody: map[string]interface{}{
				"id":   utils.RandomInt(1, 100),
				"text": utils.RandomString(6),
			},
			expectedPostRes: dto.PostResponse{
				ID:     utils.RandomInt(1, 100),
				UserID: utils.RandomInt(1, 100),
				Text:   utils.RandomString(6),
			},
			buildStubs: func(pu *mock_usecase.MockIPostUsecase, postID uint, requestBody map[string]interface{}, expectedPostRes dto.PostResponse) {
				newPost := dto.UpdatePostRequest{
					ID:   requestBody["id"].(uint),
					Text: requestBody["text"].(string),
				}
				pu.EXPECT().
					GetPostById(context.Background(), postID).
					Times(1).
					Return(expectedPostRes, nil)

				pu.EXPECT().
					UpdatePost(context.Background(), newPost).
					Times(1).
					Return(expectedPostRes, nil)
			},
			checkResponse: func(rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, rec.Code)
			},
		},
		{
			name:   "no user info",
			postID: utils.RandomInt(1, 100),
			requestBody: map[string]interface{}{
				"id":   utils.RandomInt(1, 100),
				"text": utils.RandomString(6),
			},
			expectedPostRes: dto.PostResponse{
				ID:     utils.RandomInt(1, 100),
				UserID: utils.RandomInt(1, 100),
				Text:   utils.RandomString(6),
			},
			buildStubs: func(pu *mock_usecase.MockIPostUsecase, postID uint, requestBody map[string]interface{}, expectedPostRes dto.PostResponse) {
			},
			checkResponse: func(rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, rec.Code)
			},
		},
		{
			name:   "request bind error",
			postID: utils.RandomInt(1, 100),
			requestBody: map[string]interface{}{
				"id":   "invalid id",
				"text": utils.RandomString(6),
			},
			expectedPostRes: dto.PostResponse{
				ID:     utils.RandomInt(1, 100),
				UserID: utils.RandomInt(1, 100),
				Text:   utils.RandomString(6),
			},
			buildStubs: func(pu *mock_usecase.MockIPostUsecase, postID uint, requestBody map[string]interface{}, expectedPostRes dto.PostResponse) {
			},
			checkResponse: func(rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, rec.Code)
			},
		},
		{
			name:   "internal server error",
			postID: utils.RandomInt(1, 100),
			requestBody: map[string]interface{}{
				"id":   utils.RandomInt(1, 100),
				"text": utils.RandomString(6),
			},
			expectedPostRes: dto.PostResponse{
				ID:     utils.RandomInt(1, 100),
				UserID: utils.RandomInt(1, 100),
				Text:   utils.RandomString(6),
			},
			buildStubs: func(pu *mock_usecase.MockIPostUsecase, postID uint, requestBody map[string]interface{}, expectedPostRes dto.PostResponse) {
				newPost := dto.UpdatePostRequest{
					ID:   requestBody["id"].(uint),
					Text: requestBody["text"].(string),
				}
				pu.EXPECT().
					GetPostById(context.Background(), postID).
					Times(1).
					Return(expectedPostRes, nil)

				pu.EXPECT().
					UpdatePost(context.Background(), newPost).
					Times(1).
					Return(expectedPostRes, errors.New("internal server error"))
			},
			checkResponse: func(rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, rec.Code)
			},
		},
		{
			name:   "internal server error 2",
			postID: utils.RandomInt(1, 100),
			requestBody: map[string]interface{}{
				"id":   utils.RandomInt(1, 100),
				"text": utils.RandomString(6),
			},
			expectedPostRes: dto.PostResponse{
				ID:     utils.RandomInt(1, 100),
				UserID: utils.RandomInt(1, 100),
				Text:   utils.RandomString(6),
			},
			buildStubs: func(pu *mock_usecase.MockIPostUsecase, postID uint, requestBody map[string]interface{}, expectedPostRes dto.PostResponse) {
				newPost := dto.UpdatePostRequest{
					ID:   requestBody["id"].(uint),
					Text: requestBody["text"].(string),
				}
				pu.EXPECT().
					GetPostById(context.Background(), postID).
					Times(1).
					Return(expectedPostRes, errors.New("internal server error"))

				pu.EXPECT().
					UpdatePost(context.Background(), newPost).
					Times(0)
			},
			checkResponse: func(rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, rec.Code)
			},
		},
		{
			name:   "post user id does not match token user id",
			postID: utils.RandomInt(1, 100),
			requestBody: map[string]interface{}{
				"id":   utils.RandomInt(1, 100),
				"text": utils.RandomString(6),
			},
			expectedPostRes: dto.PostResponse{
				ID:     utils.RandomInt(1, 100),
				UserID: utils.RandomInt(1, 100) + 1,
				Text:   utils.RandomString(6),
			},
			buildStubs: func(pu *mock_usecase.MockIPostUsecase, postID uint, requestBody map[string]interface{}, expectedPostRes dto.PostResponse) {
				newPost := dto.UpdatePostRequest{
					ID:   requestBody["id"].(uint),
					Text: requestBody["text"].(string),
				}
				pu.EXPECT().
					GetPostById(context.Background(), postID).
					Times(1).
					Return(expectedPostRes, nil)

				pu.EXPECT().
					UpdatePost(context.Background(), newPost).
					Times(0)

			},
			checkResponse: func(rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, rec.Code)
			},
		},
	}

	e := echo.New()
	e.Validator = &utils.CustomValidator{Validator: validator.New()}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	pu := mock_usecase.NewMockIPostUsecase(ctrl)
	pc := NewPostController(pu)

	for i := range cases {
		tc := cases[i]
		t.Run(tc.name, func(t *testing.T) {
			tc.buildStubs(pu, tc.postID, tc.requestBody, tc.expectedPostRes)

			body, err := json.Marshal(tc.requestBody)
			require.NoError(t, err)

			token, err := utils.CreateValidToken(tc.expectedPostRes.UserID)
			require.NoError(t, err)

			url := fmt.Sprintf("/posts/%d", tc.postID)
			req := httptest.NewRequest(http.MethodPut, url, bytes.NewReader(body))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %v", token))
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("postId")
			c.SetParamValues(strconv.Itoa(int(tc.postID)))

			if tc.name != "no user info" {
				user := &jwt.Token{Claims: &dto.JwtCustomClaims{ID: tc.expectedPostRes.UserID}}
				c.Set("user", user)
			}
			if tc.name == "post user id does not match token user id" {
				user := &jwt.Token{Claims: &dto.JwtCustomClaims{ID: tc.expectedPostRes.UserID + 1}}
				c.Set("user", user)
			}

			err = pc.UpdatePost(c)
			require.NoError(t, err)

			tc.checkResponse(rec)
		})
	}
}

func TestDeletePost(t *testing.T) {
	cases := []struct {
		name            string
		postID          uint
		expectedPostRes dto.PostResponse
		buildStubs      func(pu *mock_usecase.MockIPostUsecase, postID uint, expectedPostRes dto.PostResponse)
		checkResponse   func(recoder *httptest.ResponseRecorder)
	}{
		{
			name:   "valid request",
			postID: utils.RandomInt(1, 100),
			expectedPostRes: dto.PostResponse{
				ID:     utils.RandomInt(1, 100),
				UserID: utils.RandomInt(1, 100),
				Text:   utils.RandomString(6),
			},
			buildStubs: func(pu *mock_usecase.MockIPostUsecase, postID uint, expectedPostRes dto.PostResponse) {
				pu.EXPECT().
					GetPostById(context.Background(), postID).
					Times(1).
					Return(expectedPostRes, nil)

				pu.EXPECT().
					DeletePost(context.Background(), postID).
					Times(1).
					Return(nil)
			},
			checkResponse: func(rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNoContent, rec.Code)
			},
		},
		{
			name:   "no user info",
			postID: utils.RandomInt(1, 100),
			expectedPostRes: dto.PostResponse{
				ID:     utils.RandomInt(1, 100),
				UserID: utils.RandomInt(1, 100),
				Text:   utils.RandomString(6),
			},
			buildStubs: func(pu *mock_usecase.MockIPostUsecase, postID uint, expectedPostRes dto.PostResponse) {
			},
			checkResponse: func(rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, rec.Code)
			},
		},
		{
			name:   "internal server error",
			postID: utils.RandomInt(1, 100),
			expectedPostRes: dto.PostResponse{
				ID:     utils.RandomInt(1, 100),
				UserID: utils.RandomInt(1, 100),
				Text:   utils.RandomString(6),
			},
			buildStubs: func(pu *mock_usecase.MockIPostUsecase, postID uint, expectedPostRes dto.PostResponse) {
				pu.EXPECT().
					GetPostById(context.Background(), postID).
					Times(1).
					Return(expectedPostRes, nil)

				pu.EXPECT().
					DeletePost(context.Background(), postID).
					Times(1).
					Return(errors.New("internal server error"))
			},
			checkResponse: func(rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, rec.Code)
			},
		},
		{
			name:   "post user id does not match token user id",
			postID: utils.RandomInt(1, 100),
			expectedPostRes: dto.PostResponse{
				ID:     utils.RandomInt(1, 100),
				UserID: utils.RandomInt(1, 100) + 1,
				Text:   utils.RandomString(6),
			},
			buildStubs: func(pu *mock_usecase.MockIPostUsecase, postID uint, expectedPostRes dto.PostResponse) {
				pu.EXPECT().
					GetPostById(context.Background(), postID).
					Times(1).
					Return(expectedPostRes, nil)

				pu.EXPECT().
					DeletePost(context.Background(), postID).
					Times(0)
			},
			checkResponse: func(rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, rec.Code)
			},
		},
		{
			name:   "internal server error 2",
			postID: utils.RandomInt(1, 100),
			expectedPostRes: dto.PostResponse{
				ID:     utils.RandomInt(1, 100),
				UserID: utils.RandomInt(1, 100),
				Text:   utils.RandomString(6),
			},
			buildStubs: func(pu *mock_usecase.MockIPostUsecase, postID uint, expectedPostRes dto.PostResponse) {
				pu.EXPECT().
					GetPostById(context.Background(), postID).
					Times(1).
					Return(expectedPostRes, errors.New("internal server error"))

				pu.EXPECT().
					DeletePost(context.Background(), postID).
					Times(0)
			},
			checkResponse: func(rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, rec.Code)
			},
		},
	}
	e := echo.New()
	e.Validator = &utils.CustomValidator{Validator: validator.New()}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	pu := mock_usecase.NewMockIPostUsecase(ctrl)
	pc := NewPostController(pu)

	for i := range cases {
		tc := cases[i]
		t.Run(tc.name, func(t *testing.T) {
			tc.buildStubs(pu, tc.postID, tc.expectedPostRes)

			token, err := utils.CreateValidToken(tc.expectedPostRes.UserID)
			require.NoError(t, err)

			url := fmt.Sprintf("/posts/%d", tc.postID)
			req := httptest.NewRequest(http.MethodDelete, url, nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %v", token))
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("postId")
			c.SetParamValues(strconv.Itoa(int(tc.postID)))

			if tc.name != "no user info" {
				user := &jwt.Token{Claims: &dto.JwtCustomClaims{ID: tc.expectedPostRes.UserID}}
				c.Set("user", user)
			}
			if tc.name == "post user id does not match token user id" {
				user := &jwt.Token{Claims: &dto.JwtCustomClaims{ID: tc.expectedPostRes.UserID + 1}}
				c.Set("user", user)
			}

			err = pc.DeletePost(c)
			require.NoError(t, err)

			tc.checkResponse(rec)
		})
	}
}
