FROM golang:1.24.2-alpine as builder

RUN apk add --no-cache gcc musl-dev

WORKDIR /app

COPY go.mod go.sum* ./
RUN go mod download

COPY . .

RUN go build -o main .

FROM alpine:3.19

RUN apk add --no-cache \
    libc-compat \
    exfat-utils \
    fuse-exfat

WORKDIR /app

COPY --from=builder /app/main .

EXPOSE 7777

CMD ["/app/main"]