## Snap Chat Server
> This is a chat server application using golang.

Application are tested using Golang 1.22 and PostgreSQL 16.

## Clone the application 
```sh
git clone https://github.com/ezartsh/snap-chat-server.git
cd snap-chat-server
```

## Installation
### Using Docker
```sh
# No configuration needed. Just run this.
docker compose up -d
```

### Setup Database
First create / initiate the database and update the .env file for the database connection and database name.
After that migrate the database schema.
```sh
# Copy .env.example .env
cp .env.example .env
# Change the configuration based on your need.
# Migrate Database
docker run --rm -it snap-chat-server-web migrate up
```
### Manual

```sh
# After clone the repo and move inside the project folder.
go mod tidy
go mod download
# Migrate the database schema.
go run main.go migrate up
# run serve the application on port.
go run main.go serve
```
