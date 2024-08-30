# Start from golang base image
FROM golang:alpine AS builder

# Enable go modules
ENV GO111MODULE=on

# Install git. (alpine image does not have git in it)
RUN apk update && apk add --no-cache git

# Set current working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod ./
COPY go.sum ./

# Download all dependencies.
RUN go mod download

# Now, copy the source code
COPY . .

# Build the application.
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./ ./cmd/server/main.go

# Finally our multi-stage to build a small image
FROM alpine:latest

#install supervisor
RUN apk update && apk add --no-cache supervisor

#copy supervisor config
COPY etc/supervisord.conf /etc/supervisord.conf

#set work dir
WORKDIR /app

RUN mkdir /app/keys

# Copy the Pre-built binary file
COPY --from=builder /app/main .

# Change this part to using os env
COPY --from=builder /app/.env .
COPY --from=builder /app/var/cert.key ./var/

#create logs directory
RUN mkdir logs && touch logs/logger.log
#RUN chmod +x ./.env
RUN chmod 755 -R ./logs

CMD ["/usr/bin/supervisord", "-c", "/etc/supervisord.conf"]