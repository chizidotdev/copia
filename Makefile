build:
	go build -o bin/copia ./main.go

redis:
	docker run --name copia-redis -p 6389:6379 -d redis

redis-cli:
	docker exec -it copia-redis redis-cli

postgres:
	docker run --name copia -p 5433:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:12-alpine

postgres-cli:
	docker exec -it copia psql --username=root --dbname=copia

createdb:
	docker exec -it copia createdb --username=root --owner=root copia

dropdb:
	docker exec -it copia dropdb copia

test:
	go test -v ./...

server:
	air

www:
	npm --prefix client run dev

mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/chizidotdev/copia/db/sqlc Store

.PHONY: build redis redis-cli postgres postgres-cli createdb dropdb test mock
