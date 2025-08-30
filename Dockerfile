ARG GO_VERSION=1
FROM golang:${GO_VERSION}-bookworm as builder

WORKDIR /usr/src/app
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .
RUN go build -v -o /run-app ./cmd/henzworld


FROM gcr.io/distroless/base-debian12

ARG WEB_DIR=/usr/share/henzworld/
ENV STATIC_DIR=${WEB_DIR}/static/
ENV TEMPLATE_DIR=${WEB_DIR}/template/

COPY --from=builder /usr/src/app/web /usr/share/henzworld
COPY --from=builder /run-app /usr/local/bin/
CMD ["run-app"]
