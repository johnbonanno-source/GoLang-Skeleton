# CarTrader

## General

Parts of the system are configured using environment variables. For local development, those are not secret and local development values are currently in the `.env.example` file.

To get started, copy the `.env.example` to `.env`

```
cp .env.example .env
```

## Backend

### Prerequisites

1. golang - 1.24
2. docker
3. goose (migration tool) -https://github.com/pressly/goose

### Development

To start the database and run the migrations

```
docker compose up

goose up
```

To start the backend server

```
go run ./backend/cmd
```

## Frontend

### Prerequisites

1. node - 22.14 (lts)
2. pnpm - https://pnpm.io/installation

### Development

To get started with the frontend, do the following:

```
pnpm install
pnpm dev
```
