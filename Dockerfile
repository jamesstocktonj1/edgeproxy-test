# syntax=docker/dockerfile:1

# build container
FROM golang:1.22 AS build

WORKDIR /app

COPY . ./
RUN go mod download

ARG CGO_ENABLED=0
RUN go build -o ./server ./cmd/server

# run container
FROM envoyproxy/envoy:dev

WORKDIR /app

COPY --from=build /app/server ./server

COPY proxy/edgeproxy.yaml /etc/envoy.yaml

COPY start.sh ./start.sh

EXPOSE 8080
EXPOSE 9090

CMD ["sh", "start.sh"]