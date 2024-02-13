package main

import (
	"database/sql"
	"log"

	"github.com/PenginAction/go-BulletinBoard/config"
	"github.com/PenginAction/go-BulletinBoard/controller"
	db "github.com/PenginAction/go-BulletinBoard/db/sqlc"
	"github.com/PenginAction/go-BulletinBoard/router"
	"github.com/PenginAction/go-BulletinBoard/usecase"

	_ "github.com/lib/pq"
)

func main() {
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}
	conn, err := sql.Open(cfg.DBDriver, cfg.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	store := db.NewStore(conn)

	userUsecase := usecase.NewUserUsecase(store)
	postUsecase := usecase.NewPostUsecase(store)
	userController := controller.NewUserController(userUsecase)
	postController := controller.NewPostController(postUsecase)

	e := router.NewRouter(userController, postController, cfg)
	e.Logger.Fatal(e.Start(":8080"))
}
