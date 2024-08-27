postgres:
	docker run --name go-boiler-postgres -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=boiler -d postgres:latest

createdb:
	docker exec -it go-boiler-postgres createdb --username=root --owner=root simple_bank

dropdb:
	docker exec -it go-boiler-postgres dropdb  --username=root  simple_bank

migrateup:
	migrate -path db/migration -database "postgresql://root:boiler@localhost:5432/simple_bank?sslmode=disable" -verbose up
	
migratedown:
	migrate -path db/migration -database "postgresql://root:boiler@localhost:5432/simple_bank?sslmode=disable" -verbose down 

sqlc:
	sqlc generate 

test:
	go test -v -cover ./...

server:
	go run main.go 

mock:
	mockgen -package mockdb  -destination db/mock/store.go github.com/hex-aragon/go-backend-boilerplate/db/sqlc Store 

.PHONY: postgres createdb dropdb migrateup migratedown sqlc test server mock 