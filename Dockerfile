FROM golang:alpine AS builder

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /dist

COPY go.mod .
COPY go.sum .
RUN go mod download
RUN go mod verify
COPY . .
RUN go build -o main cmd/*/main.go

FROM gcr.io/distroless/static

WORKDIR /app

COPY --from=builder /dist/main .
COPY db/migrations db/migrations
EXPOSE 8080

ENTRYPOINT [ "/app/main" ]
