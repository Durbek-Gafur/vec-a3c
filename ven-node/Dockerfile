FROM golang:1.20 

ENV PROJECT_DIR=/app \
    GO111MODULE=on \
    CGO_ENABLED=0 \
    GOPROXY=https://proxy.golang.org,direct \
    GONOPROXY=none \
    GOSUMDB=sum.golang.org \
    GONOSUMDB=none

WORKDIR /app

COPY ./ /app

RUN go mod tidy
RUN go build -o /build/app ./cmd/vec_worker

ENTRYPOINT ["/build/app"]
