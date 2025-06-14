postgres:
	docker run --name postgres12 --rm -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=123456 -d postgres:12-alpine

createdb:
	docker exec postgres12 createdb --user=root --owner=root bank

dropdb:
	docker exec postgres12 dropdb --user=root bank

migrate-up:
	migrate -path db/migration -database "postgres://root:123456@localhost/bank?sslmode=disable" -verbose up

migrate-down:
	migrate -path db/migration -database "postgres://root:123456@localhost/bank?sslmode=disable" -verbose down

dbdocs:
	dbdocs build doc/db.dbml --project bank

dbschema:
	dbml2sql --postgres -o doc/schema.sql doc/db.dbml

sqlc:
	sqlc generate

server:
	go run main.go

mock:
	mockgen -destination db/mock/store.go -package mockdb github.com/Nickeymaths/bank/db/sqlc Store

test:
	go clean -testcache && go test -v -cover ./...

proto:
	rm -r pb/*.go
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
    --go-grpc_out=pb --go-grpc_opt=paths=source_relative \
	--grpc-gateway_out=pb --grpc-gateway_opt=paths=source_relative \
    proto/*.proto

evans:
	evans --host localhost --port 5000 -r repl

.PHONY: postgres mysql-up createdb dropdb migrate-up migrate-down dbdocs dbschema sqlc server mockdb test proto evans
