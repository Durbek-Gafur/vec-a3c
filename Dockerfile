FROM golang:latest

ENV PROJECT_DIR=/app \
    GO111MODULE=on \
    CGO_ENABLED=0 \
    GOPROXY=https://proxy.golang.org,direct \
    GONOPROXY=none \
    GOSUMDB=sum.golang.org \
    GONOSUMDB=none

WORKDIR /app
RUN mkdir "/build"

COPY ./ /app


RUN go install -mod=mod github.com/githubnemo/CompileDaemon
RUN go mod tidy


ENTRYPOINT CompileDaemon -build="go build -o /build/app ./cmd/vec_worker" -command="/build/app"



