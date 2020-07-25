FROM jordimartin/mmock:v3.0.0 as release-mmock
FROM golang:1-alpine as build-spawner

RUN set -ex \
  && apk add --no-cache -q --no-progress build-base git linux-headers

WORKDIR /go/app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build -o spawner .

FROM alpine:3

RUN set -ex \
  && apk add --no-cache -q --no-progress bash curl ca-certificates tzdata \
  && update-ca-certificates \
  && rm -rf /var/cache/apk/* /tmp/* \
  && mkdir -p /config \
  && mkdir -p /tls

VOLUME /config

COPY --from=release-mmock /tls/server.crt /tls/server.crt
COPY --from=release-mmock /tls/server.key /tls/server.key
COPY --from=release-mmock /usr/local/bin/mmock /usr/local/bin/mmock
COPY --from=build-spawner /go/app/spawner /usr/local/bin/spawner

EXPOSE 5000 8082 8083 8084

CMD ["spawner"]
