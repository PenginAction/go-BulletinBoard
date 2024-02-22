package router

import (
	"github.com/PenginAction/go-BulletinBoard/config"
	"github.com/PenginAction/go-BulletinBoard/controller"
	"github.com/PenginAction/go-BulletinBoard/dto"
	"github.com/PenginAction/go-BulletinBoard/utils"
	"github.com/go-playground/validator"
	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func NewRouter(uc controller.IUserController, pc controller.IPostController, cfg config.Config) *echo.Echo {
	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowCredentials: true,
	}))
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "time=${time_rfc3339_nano}, method=${method}, uri=${uri}, status=${status}\n",
	}))
	e.Validator = &utils.CustomValidator{Validator: validator.New()}

	e.POST("/signup", uc.Signup)
	e.POST("/login", uc.Login)
	// r.POST("/logout", uc.Logout)

	p := e.Group("/posts")

	config := echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(dto.JwtCustomClaims)
		},
		SigningKey: []byte(cfg.SECRET),
	}

	p.Use(echojwt.WithConfig(config))
	p.GET("", pc.GetAllPosts)
	p.GET("/:postId", pc.GetPostById)
	p.POST("", pc.CreatePost)
	p.PUT("/:postId", pc.UpdatePost)
	p.DELETE("/:postId", pc.DeletePost)

	return e
}
