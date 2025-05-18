.PHONY: postgres-up mysql-up postgres-down createdb dropdb migrate-up migrate-down sqlc server mockdb test

postgres-up:
	docker run --name postgres12 --network bank-network --rm -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=123456 -d postgres:12-alpine

postgres-down:
	docker stop postgres12

createdb:
	docker exec postgres12 createdb --user=root --owner=root bank

dropdb:
	docker exec postgres12 dropdb --user=root bank
run:
	docker run -d -p 4000:4000 --name bank --network bank-network --rm -e GIN_MOD=release -e DB_SOURCE='postgres://root:123456@postgres12/bank?sslmode=disable' bank:v1

migrate-up:
	migrate -path db/migration -database "postgres://root:123456@postgres12/bank?sslmode=disable" -verbose up

migrate-down:
	migrate -path db/migration -database "postgres://root:123456@postgres12/bank?sslmode=disable" -verbose down

migrate-up1:
	migrate -path db/migration -database "postgres://root:123456@postgres12/bank?sslmode=disable" -verbose up 1

migrate-down1:
	migrate -path db/migration -database "postgres://root:123456@postgres12/bank?sslmode=disable" -verbose down 1

sqlc:
	sqlc generate

server:
	go run main.go

mockdb:
	mockgen -destination db/mock/store.go -package mockdb github.com/Nickeymaths/bank/db/sqlc Store

test:
	go clean -testcache && go test -v -cover ./...
