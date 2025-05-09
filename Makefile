.PHONY: postgres-up mysql-up postgres-down createdb dropdb migrate-up migrate-down sqlc server mock test

postgres-up:
	docker run --name postgres12 --rm -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=123456 -d postgres:12-alpine

postgres-down:
	docker stop postgres12

createdb:
	docker exec postgres12 createdb --user=root --owner=root bank

dropdb:
	docker exec postgres12 dropdb --user=root bank

migrate-up:
	migrate -path db/migration -database "postgres://root:123456@localhost/bank?sslmode=disable" -verbose up

migrate-down:
	migrate -path db/migration -database "postgres://root:123456@localhost/bank?sslmode=disable" -verbose down

migrate-up1:
	migrate -path db/migration -database "postgres://root:123456@localhost/bank?sslmode=disable" -verbose up 1

migrate-down1:
	migrate -path db/migration -database "postgres://root:123456@localhost/bank?sslmode=disable" -verbose down 1

sqlc:
	sqlc generate

server:
	go run main.go

mock:
	mockgen -destination db/mock/store.go -package mockdb github.com/Nickeymaths/bank/db/sqlc Store

test:
	go clean -testcache && go test -v -cover ./...
