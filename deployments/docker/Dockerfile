FROM golang:1.16

ENV build_env=docker

WORKDIR /app

COPY ./ /app

RUN go mod download

RUN go get github.com/githubnemo/CompileDaemon

EXPOSE 8080

ENTRYPOINT CompileDaemon --build="go build cmd/product_svc/main.go" --command=./main