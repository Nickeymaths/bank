.PHONY: postgres-up postgres-down createdb dropdb migrate-up migrate-down sqlc test

postgres-up:
	docker run --name postgres-17 --rm -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=123456 -d postgres:17-alpine

postgres-down:
	docker stop postgres-17

createdb:
	docker exec postgres-17 createdb --user=root --owner=root bank

dropdb:
	docker exec postgres-17 dropdb --user=root bank

migrate-up:
	migrate -path db/migration -database "postgres://root:123456@localhost/bank?sslmode=disable" up

migrate-down:
	migrate -path db/migration -database postgres://root:123456@localhost/bank?sslmode=disable down

sqlc:
	sqlc generate

test:
	go test -v -cover ./...