# Goecommerce API

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
    - ensure configure your environment variables in .env file (copy from .env.example file).
    - docker build -t gocommerce-api .
    - docker run -it -p 3001:3001 --rm --name gocommerce-api gocommerce-api

## Endpoints

See Doc [API Endpoint](./ENDPOINT.md)

## Test

1. E2E Test, you can test the API endpoints using via this file [Http Test File](./test/http/http_test.http)
2. `go test ./..` still under development