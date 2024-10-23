postgres:
	docker run --name go-boiler-postgres --network bank-network -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=boiler -d postgres:latest

createdb:
	docker exec -it go-boiler-postgres createdb --username=root --owner=root simple_bank

dropdb:
	docker exec -it go-boiler-postgres dropdb  --username=root  simple_bank

migrateup:
	migrate -path db/migration -database "postgresql://root:boiler@localhost:5432/simple_bank?sslmode=disable" -verbose up

migrateup1:
	migrate -path db/migration -database "postgresql://root:boiler@localhost:5432/simple_bank?sslmode=disable" -verbose up 1
		
migratedown:
	migrate -path db/migration -database "postgresql://root:boiler@localhost:5432/simple_bank?sslmode=disable" -verbose down 
	
migratedown1:
	migrate -path db/migration -database "postgresql://root:boiler@localhost:5432/simple_bank?sslmode=disable" -verbose down 1


sqlc:
	sqlc generate 

test:
	go test -v -cover ./...

server:
	go run main.go 

mock:
	mockgen -package mockdb  -destination db/mock/store.go github.com/hex-aragon/go-backend-boilerplate/db/sqlc Store 

docker-build:
	docker build -t simplebank:latest .

docker-run-release:
	docker run --name simplebank -p 8080:8080 --network bank-network -e GIN_MODE=release -e DB_SOURCE="postgresql://root:boiler@go-boiler-postgres:5432/simple_bank?sslmode=disable" simplebank:latest

.PHONY: postgres createdb dropdb migrateup migrateup1 migratedown migratedown1 sqlc test server mock docker-build docker-run