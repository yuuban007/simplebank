postgres:
	docker run --name postgres12 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=gcc123456 -d postgres:12-alpine
createdb:
	docker exec -it postgres12 createdb --username=root --owner=root simple_bank
dropdb:
	docker exec -it postgres12 dropdb simple_bank
migrationup:
	migrate -path db/migration -database "postgresql://root:gcc123456@localhost:5432/simple_bank?sslmode=disable" -verbose up
migrationdown:
	migrate -path db/migration -database "postgresql://root:gcc123456@localhost:5432/simple_bank?sslmode=disable" -verbose down
sqlc:
	sqlc generate
test:
	go test -v -cover ./...
server:
	go run main.go
.PHONY: postgres,createdb,dropdb,sqlc,test,server