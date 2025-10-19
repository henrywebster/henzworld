ARG GO_VERSION=1
FROM golang:${GO_VERSION}-bookworm AS builder

WORKDIR /usr/src/app
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .
RUN go build -v -o /run-app ./cmd/henzworld

FROM debian:bookworm-slim

RUN apt-get update && apt-get install -y \
    nginx \
    supervisor \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/* \
    && update-ca-certificates

ARG WEB_DIR=/usr/share/henzworld/
ENV STATIC_DIR=${WEB_DIR}/static/
ENV TEMPLATE_DIR=${WEB_DIR}/template/

ARG DB_FILE=lm-data.db
ENV DATABASE_LOCAL_FILE=/var/lib/app/${DB_FILE}

COPY --from=builder /usr/src/app/web /usr/share/henzworld
COPY --from=builder /run-app /usr/local/bin/

COPY ${DB_FILE} /var/lib/app/${DB_FILE}
COPY nginx.conf /etc/nginx/nginx.conf
COPY supervisord.conf /etc/supervisor/conf.d/supervisord.conf

RUN mkdir -p /var/cache/nginx

EXPOSE 8080
CMD ["/usr/bin/supervisord", "-c", "/etc/supervisor/conf.d/supervisord.conf"]
