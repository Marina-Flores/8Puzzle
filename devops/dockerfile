FROM golang:1.23 AS builder
WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .
RUN go build -o 8Puzzle .

FROM debian:bullseye-slim
COPY --from=builder /app/8Puzzle /usr/local/bin/8Puzzle
CMD ["8Puzzle"]
