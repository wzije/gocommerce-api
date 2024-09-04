# Gocommerce API

Sample Ecommerce application using Go

## Requirement

- golang v1.22
- postgresql v13
- docker

## How To Run

1. Locally

   - Clone the repository.
   - Install dependencies using `go mod download`.
   - Configure your environment variables in .env file (copy from .env.example file).
   - Run the application using `go run ./cmd/server/main.go`.

2. docker

   Ensure configure your environment variables in .env file (copy from .env.example file).
   Run this command

   ```bash
   $ docker build -t gocommerce-api .
   $ docker run -it -p 3001:3001 --rm --name gocommerce-api gocommerce-api
   ```

## Migration

1. Create database, you can manually create the database and table or Run
   this [Schema](./db/migration/000001_init_schema.up.sql) to your
   database command.
2. Execute the dummy data / seeder using `go run ./cmd/seeds/seed.go`

## Endpoints

See Doc [API Endpoint](./ENDPOINT.md)

## Test

1. E2E Test, you can test the API endpoints

   - Ensure the seed and application was running
     ```bash
     $ go run cmd/seeds/seed.go
     $ go run cmd/server/main.go
     ```
   - Open this file [Http Test File](./test/http/http_test.http) and run All Requests (click run all request button)

     > _Note: Only support on goland or jetbrain plathform_

2. `go test $(go list ./...) -v` -> **Under Development**
