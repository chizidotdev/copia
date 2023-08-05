postgres:
	sudo docker run --name copia -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:12-alpine

createdb:
	sudo docker exec -it copia createdb --username=root --owner=root copia

dropdb:
	sudo docker exec -it copia dropdb copia

test:
	go test -v ./...

server:
	go run main.go

www:
	npm --prefix client run dev

mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/chizidotdev/copia/db/sqlc Store

.PHONY: postgres createdb dropdb test mock
