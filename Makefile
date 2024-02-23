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
