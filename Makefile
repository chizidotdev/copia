postgres:
	docker run --name postgres12 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:12-alpine

postgres-cli:
	docker exec -it postgres12 psql --username=root --dbname=am-shop

createdb:
	docker exec -it postgres12 createdb --username=root --owner=root am-shop

dropdb:
	docker exec -it postgres12 dropdb am-shop

sqlc:
	sqlc generate

redis:
	docker run --name am-shop-redis -p 6389:6379 -d redis

redis-cli:
	docker exec -it am-shop-redis redis-cli

server:
	air

migrate-up:
	go run cmd/migrations/init/init.sql.go up

migrate-down:
	go run cmd/migrations/init/init.sql.go down


.PHONY: postgres postgres-cli createdb dropdb sqlc redis redis-cli server
