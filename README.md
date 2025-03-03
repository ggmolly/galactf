# setup

## dependencies

1. go 1.24.0
2. yarn 1.22.x

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