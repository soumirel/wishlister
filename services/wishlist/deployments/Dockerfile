FROM golang:1.25-alpine AS builder
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 \
    GOOS=linux \
    go build -ldflags="-w -s" -o /app/wishlister /src/cmd/wishlister/main.go 

FROM alpine:3.23
RUN addgroup -g 1001 -S app && \
    adduser -u 1001 -S -D -H -G app app

WORKDIR /app

COPY --from=builder --chown=app:app /app/wishlister /app/wishlister

COPY --from=builder --chown=app:app /src/db/init/init.sql /app/scripts/
COPY --from=builder --chown=app:app /src/db/refresh_test_data.sql /app/scripts/

COPY --from=builder --chown=app:app /src/deployments/config.compose.yaml /app/config.yaml

USER 1001

# http
EXPOSE 8080
# grpc
EXPOSE 8081

CMD ["/app/wishlister"]