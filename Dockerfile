# run test and build stage
FROM golang:1.18-alpine3.15 AS builder
WORKDIR /app
COPY go.mod /app/go.mod
COPY go.sum /app/go.sum
RUN go mod download
#RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /main
COPY . .
# disable CGO to fix missing gcc: `CGO_ENABLED=0`
#RUN CGO_ENABLED=0 go test ./...
RUN go build -o main main.go

# run stage
FROM alpine:3.15
WORKDIR /app
COPY --from=builder /app/main .
COPY .env .
COPY /db/migrations /app/db/migrations

EXPOSE 8000
CMD [ "/app/main" ]
