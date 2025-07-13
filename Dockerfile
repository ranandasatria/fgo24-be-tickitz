ARG GO_VERSION=1.24.4
FROM golang:${GO_VERSION}-alpine AS builder

WORKDIR /buildapp

COPY go.mod go.sum ./
RUN go mod download


COPY . .
RUN CGO_ENABLED=0 go build -o goapp .
RUN CGO_ENABLED=0 go build -o seed_tmdb ./cmd/seed_tmdb


FROM alpine:3.22

WORKDIR /app
COPY --from=builder /buildapp/goapp /app/goapp
COPY --from=builder /buildapp/seed_tmdb /app/seed_tmdb

ENTRYPOINT ["/app/goapp"]
