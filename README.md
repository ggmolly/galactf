# setup

## dependencies

1. go 1.24.0
2. yarn 1.22.x
3. postgres 17.x (see [docker-compose.yml](docker-compose.yml))
4. redis 7.4.x (see [docker-compose.yml](docker-compose.yml))

## websocket dev dependencies

1. go protobuf compiler 1.36.5 (`go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.36.5`)
2. proto 29.3-1
3. ts-proto ^2.6.1 (in `front/node_modules`)

## env

1. `cd back`
2. `cp .env.sample .env`

## front

1. `cd front`
2. `yarn`
3. `yarn dev`

> [!NOTE]
> Will run a development React server on port 5173

## back

1. `cd back`
2. `go mod tidy`
3. get [air](https://github.com/air-verse/air) `go install github.com/air-verse/air@latest`
4. `air`

> [!NOTE]
> Will run a development server on port 7777

# seeding

1. `cd back`
2. `go run main.go seed`
