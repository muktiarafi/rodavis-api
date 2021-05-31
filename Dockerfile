FROM golang:alpine AS builder

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /dist
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN go build -o main cmd/rodavis/main.go

FROM alpine

WORKDIR /app
COPY --from=builder /dist/main .
COPY db/migrations db/migrations
COPY key.json .
EXPOSE 8080

ENTRYPOINT [ "/app/main" ]
