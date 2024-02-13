postgres:
	docker run --name postgres16 --network bulletin-board-network -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:16-alpine

createdb:
	docker exec -it postgres16 createdb --username=root --owner=root bulletin_board

dropdb:
	docker exec -it postgres16 dropdb bulletin_board
	
migrateup:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/bulletin_board?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/bulletin_board?sslmode=disable" -verbose down

sqlc:
	sqlc generate

test: 
	go test -v -cover ./...
	
start:
	go run main.go

fmt:
	go fmt ./...

mockdb:
	mockgen -package mockdb -destination db/mock/store.go github.com/PenginAction/go-BulletinBoard/db/sqlc Store 

mockuser:
	mockgen -source usecase/user_usecase.go -destination usecase/mock/UserUsecase.go

mockpost:
	mockgen -source usecase/post_usecase.go -destination usecase/mock/PostUsecase.go

mockimage:
	mockgen -source usecase/image_usecase.go -destination usecase/mock/ImageUsecase.go

.PHONY: postgres createdb dropdb migrateup migratedown sqlc test start fmt mockdb mockuser mockpost mockimage
