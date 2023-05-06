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

RUN go install -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
RUN go install github.com/golang/mock/mockgen@v1.6.0
RUN go install -mod=mod github.com/githubnemo/CompileDaemon
RUN go mod tidy

ENTRYPOINT CompileDaemon -polling=true -build="go build -o /build/app ./cmd/vec_worker" -command="/build/app"



