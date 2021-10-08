FROM golang:1.16 AS build

ENV build_env=docker

WORKDIR /app

COPY ./ /app
RUN go mod download

RUN go build -o /product-svc cmd/product_svc/main.go

FROM gcr.io/distroless/base-debian10

WORKDIR /

COPY --from=build /product-svc /product-svc
COPY --from=build /app/config /config

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/product-svc"]