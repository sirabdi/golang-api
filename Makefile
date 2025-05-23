DB_PASSWORD?= admin
DB_NAME ?= simple_bank
DB_USER ?= postgres
DB_HOST ?= localhost
DB_PORT ?= 5432

createdb:
	@powershell -Command "$$env:PGPASSWORD='$(DB_PASSWORD)'; createdb -h $(DB_HOST) -p $(DB_PORT) -U $(DB_USER) $(DB_NAME)"

dropdb:
	@powershell -Command "$$env:PGPASSWORD='$(DB_PASSWORD)'; dropdb -h $(DB_HOST) -p $(DB_PORT) -U $(DB_USER) $(DB_NAME)"

migrateup:
	migrate -path db/migration -database "postgresql://postgres:admin@localhost:5432/simple_bank?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://postgres:admin@localhost:5432/simple_bank?sslmode=disable" -verbose down

sqlc:
	sqlc generate

.PHONY: createdb dropdb sqlc migrateup migratedown
