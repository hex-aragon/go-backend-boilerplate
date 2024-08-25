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

.PHONY: postgres createdb dropdb migrateup migratedown