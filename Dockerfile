FROM golang:latest

ENV PROJECT_DIR=/app \
    GO111MODULE=on \
    CGO_ENABLED=0

WORKDIR /app
RUN mkdir "/build"

COPY ./ /app

RUN go mod tidy
RUN go install -mod=mod github.com/githubnemo/CompileDaemon
ENTRYPOINT CompileDaemon -polling=true -build="go build -o /build/app ./cmd/vec_worker" -command="/build/app"



