DB_URL=postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable

postgres:
	docker run --name postgres12 --network bank-network -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:12-alpine

createdb:
	docker exec -it postgres12 createdb --username=root --owner=root simple_bank

migrateup:
	migrate -path db/migration -database "$(DB_URL)" -verbose up

migrateup1:
	migrate -path db/migration -database "$(DB_URL)" -verbose up 1

migratedown:
	migrate -path db/migration -database "$(DB_URL)" -verbose down

migratedown1:
	migrate -path db/migration -database "$(DB_URL)" -verbose down 1

dropdb:
	docker exec -it postgres12 dropdb simple_bank

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

start:
	docker compose up -d

stop:
	docker compose down

dev:
	docker compose -f docker-compose.dev.yml up -d && make server

mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/shouta0715/simple-bank/db/sqlc Store

db_docs:
	dbdocs build doc/db.dbml

db_schme:
	dbml2sql --postgres -o doc/shema.sql doc/db.dbml

proto:
rm -f pb/*.go
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
    --go-grpc_out=pb --go-grpc_opt=paths=source_relative \
    proto/*.proto

.PHONY: postgres createdb migrateup migratedown dropdb sqlc test server mock migratedown1 migrateup1 start stop dev db_docs db_schme proto
