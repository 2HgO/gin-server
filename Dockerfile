FROM golang:1.13-stretch AS builder

RUN mkdir -p /go/src/gin-server

WORKDIR /go/src/gin-server

COPY . .

RUN GIT_TERMINAL_PROMPT=1 \
    GOARCH=amd64 \
    GOOS=linux \
    CGO_ENABLED=0 \
    go build -v --installsuffix cgo --ldflags="-s" -o gin-server

FROM alpine:latest

RUN mkdir -p /app/gin-server
COPY --from=builder /go/src/gin-server/gin-server /app/gin-server/

WORKDIR /app/gin-server

RUN apk add --no-cache tzdata

CMD ["./gin-server"]