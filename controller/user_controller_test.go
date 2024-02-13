package controller

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/PenginAction/go-BulletinBoard/dto"
	mock_usecase "github.com/PenginAction/go-BulletinBoard/usecase/mock"
	"github.com/PenginAction/go-BulletinBoard/utils"
	"github.com/go-playground/validator"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
)

func TestSignUp(t *testing.T) {
	cases := []struct {
		name          string
		requestBody   map[string]interface{}
		buildStubs    func(uu *mock_usecase.MockIUserUsecase)
		checkResponse func(recoder *httptest.ResponseRecorder)
	}{
		{
			name: "valid request",
			requestBody: map[string]interface{}{
				"user_str_id": utils.RandomUserStrID(),
				"email":       utils.RandomEmail(),
				"password":    utils.RandomString(6),
			},
			buildStubs: func(uu *mock_usecase.MockIUserUsecase) {
				uu.EXPECT().
					SignUp(context.Background(), gomock.Any()).
					Times(1).
					Return(dto.CreateUserResponse{}, nil)
			},
			checkResponse: func(rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusCreated, rec.Code)
			},
		},
		{
			name: "invalid request body",
			requestBody: map[string]interface{}{
				"user_str_id": utils.RandomUserStrID(),
				"email":       utils.RandomEmail(),
			},
			buildStubs: func(uu *mock_usecase.MockIUserUsecase) {

			},
			checkResponse: func(rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, rec.Code)
			},
		},
		{
			name: "invalid user_str_id 1",
			requestBody: map[string]interface{}{
				"user_str_id": 123456,
				"email":       utils.RandomEmail(),
				"password":    utils.RandomString(6),
			},
			buildStubs: func(uu *mock_usecase.MockIUserUsecase) {
				uu.EXPECT().
					SignUp(context.Background(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, rec.Code)
			},
		},
		{
			name: "invalid user_str_id 2",
			requestBody: map[string]interface{}{
				"user_str_id": "invalid_user_str_id&",
				"email":       utils.RandomEmail(),
				"password":    utils.RandomString(6),
			},
			buildStubs: func(uu *mock_usecase.MockIUserUsecase) {
				uu.EXPECT().
					SignUp(context.Background(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, rec.Code)
			},
		},
		{
			name: "invalid email 1",
			requestBody: map[string]interface{}{
				"user_str_id": utils.RandomUserStrID(),
				"email":       "invalid_email",
				"password":    utils.RandomString(6),
			},
			buildStubs: func(uu *mock_usecase.MockIUserUsecase) {
				uu.EXPECT().
					SignUp(context.Background(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, rec.Code)
			},
		},
		{
			name: "invalid email 2",
			requestBody: map[string]interface{}{
				"user_str_id": utils.RandomUserStrID(),
				"email":       123456,
				"password":    utils.RandomString(6),
			},
			buildStubs: func(uu *mock_usecase.MockIUserUsecase) {
				uu.EXPECT().
					SignUp(context.Background(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, rec.Code)
			},
		},
		{
			name: "short password",
			requestBody: map[string]interface{}{
				"user_str_id": utils.RandomUserStrID(),
				"email":       utils.RandomEmail(),
				"password":    "abc",
			},
			buildStubs: func(uu *mock_usecase.MockIUserUsecase) {
				uu.EXPECT().
					SignUp(context.Background(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, rec.Code)
			},
		},
		{
			name: "invalid password",
			requestBody: map[string]interface{}{
				"user_str_id": utils.RandomUserStrID(),
				"email":       utils.RandomEmail(),
				"password":    123456,
			},
			buildStubs: func(uu *mock_usecase.MockIUserUsecase) {
				uu.EXPECT().
					SignUp(context.Background(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, rec.Code)
			},
		},
		{
			name: "internal server error",
			requestBody: map[string]interface{}{
				"user_str_id": utils.RandomUserStrID(),
				"email":       utils.RandomEmail(),
				"password":    utils.RandomString(6),
			},
			buildStubs: func(uu *mock_usecase.MockIUserUsecase) {
				uu.EXPECT().
					SignUp(context.Background(), gomock.Any()).
					Times(1).
					Return(dto.CreateUserResponse{}, errors.New("internal server error"))
			},
			checkResponse: func(rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, rec.Code)
			},
		},
		{
			name: "invalid request body",
			requestBody: map[string]interface{}{
				"user_str_id": utils.RandomUserStrID(),
				"email":       utils.RandomEmail(),
			},
			buildStubs: func(uu *mock_usecase.MockIUserUsecase) {
				uu.EXPECT().
					SignUp(context.Background(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, rec.Code)
			},
		},
	}

	e := echo.New()
	e.Validator = &utils.CustomValidator{Validator: validator.New()}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	uu := mock_usecase.NewMockIUserUsecase(ctrl)
	uc := NewUserController(uu)

	for i := range cases {
		tc := cases[i]
		t.Run(tc.name, func(t *testing.T) {
			tc.buildStubs(uu)

			body, err := json.Marshal(tc.requestBody)
			require.NoError(t, err)

			url := "/signup"
			req := httptest.NewRequest(http.MethodPost, url, bytes.NewReader(body))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			err = uc.Signup(c)
			require.NoError(t, err)

			tc.checkResponse(rec)
		})

	}
}

func TestLogin(t *testing.T) {
	cases := []struct {
		name          string
		requestBody   map[string]interface{}
		buildStubs    func(uu *mock_usecase.MockIUserUsecase)
		checkResponse func(recoder *httptest.ResponseRecorder)
	}{
		{
			name: "valid request",
			requestBody: map[string]interface{}{
				"email":    utils.RandomEmail(),
				"password": utils.RandomString(6),
			},
			buildStubs: func(uu *mock_usecase.MockIUserUsecase) {
				uu.EXPECT().
					Login(context.Background(), gomock.Any()).
					Times(1).
					Return("test_token", nil)
			},
			checkResponse: func(rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, rec.Code)
			},
		},
		{
			name: "invalid email 1",
			requestBody: map[string]interface{}{
				"email":    "invalid_email",
				"password": utils.RandomString(6),
			},
			buildStubs: func(uu *mock_usecase.MockIUserUsecase) {
				uu.EXPECT().
					Login(context.Background(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, rec.Code)
			},
		},
		{
			name: "invalid email 2",
			requestBody: map[string]interface{}{
				"email":    123456,
				"password": utils.RandomString(6),
			},
			buildStubs: func(uu *mock_usecase.MockIUserUsecase) {
				uu.EXPECT().
					Login(context.Background(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, rec.Code)
			},
		},
		{
			name: "short password",
			requestBody: map[string]interface{}{
				"email":    utils.RandomEmail(),
				"password": "abc",
			},
			buildStubs: func(uu *mock_usecase.MockIUserUsecase) {
				uu.EXPECT().
					Login(context.Background(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, rec.Code)
			},
		},
		{
			name: "short password",
			requestBody: map[string]interface{}{
				"email":    utils.RandomEmail(),
				"password": "abc",
			},
			buildStubs: func(uu *mock_usecase.MockIUserUsecase) {
				uu.EXPECT().
					Login(context.Background(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, rec.Code)
			},
		},
		{
			name: "internal server error",
			requestBody: map[string]interface{}{
				"email":    utils.RandomEmail(),
				"password": utils.RandomString(6),
			},
			buildStubs: func(uu *mock_usecase.MockIUserUsecase) {
				uu.EXPECT().
					Login(context.Background(), gomock.Any()).
					Times(1).
					Return("", errors.New("internal server error"))
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

	uu := mock_usecase.NewMockIUserUsecase(ctrl)
	uc := NewUserController(uu)

	for i := range cases {
		tc := cases[i]
		t.Run(tc.name, func(t *testing.T) {
			tc.buildStubs(uu)
			body, err := json.Marshal(tc.requestBody)
			require.NoError(t, err)

			url := "/signup"
			req := httptest.NewRequest(http.MethodPost, url, bytes.NewReader(body))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			err = uc.Login(c)
			require.NoError(t, err)

			tc.checkResponse(rec)
		})
	}

}
