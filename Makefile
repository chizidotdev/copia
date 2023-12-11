build:
	go build -o bin/copia ./main.go

redis:
	sudo docker run --name copia-redis -p 6389:6379 -d redis

redis-cli:
	sudo docker exec -it copia-redis redis-cli

postgres:
	sudo docker run --name copia -p 5433:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:12-alpine

postgres-cli:
	sudo docker exec -it copia psql --username=root --dbname=copia

createdb:
	sudo docker exec -it copia createdb --username=root --owner=root copia

dropdb:
	sudo docker exec -it copia dropdb copia

test:
	go test -v ./...

server:
	air

www:
	npm --prefix client run dev

mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/chizidotdev/copia/db/sqlc Store

.PHONY: build redis redis-cli postgres postgres-cli createdb dropdb test mock
