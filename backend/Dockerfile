ARG GO_VERSION=1.24
FROM golang:${GO_VERSION} AS builder
WORKDIR /app
COPY . .
RUN go mod download && CGO_ENABLED=0 go build -o server ./cmd/server

FROM gcr.io/distroless/base-debian12
WORKDIR /app
COPY --from=builder /app/server /app/server
COPY --from=builder /app/db /app/db
COPY --from=builder /app/internal /app/internal
ENV APP_PORT=8280
EXPOSE 8280
CMD ["/app/server"]
