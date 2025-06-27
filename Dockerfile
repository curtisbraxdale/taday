ARG GO_VERSION=1
FROM golang:${GO_VERSION}-bookworm AS builder

WORKDIR /usr/src/app
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .
RUN go build -v -o /taday-api ./cmd/api


FROM debian:bookworm

COPY --from=builder /taday-api /usr/local/bin/
CMD ["taday-api"]
