FROM golang:1.24.2-alpine as builder

RUN apk add --no-cache gcc musl-dev

WORKDIR /app

COPY go.mod go.sum* ./
RUN go mod download

COPY go.mod go.mod
COPY go.sum go.sum
COPY main.go main.go
COPY assets assets
COPY cache cache
COPY dto dto
COPY factories factories
COPY middlewares middlewares
COPY orm orm
COPY proto proto
COPY routes routes
COPY types types
COPY utils utils
COPY views views

RUN go build -o main .

FROM alpine:3.19

RUN apk add --no-cache \
    libc6-compat \
    exfat-utils \
    fuse-exfat \
    exfatprogs \
    tzdata

WORKDIR /app

COPY .env .env
COPY --from=builder /app/main .

EXPOSE 7777

CMD ["/app/main"]