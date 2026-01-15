FROM golang:1.23-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -tags netgo -ldflags "-s -w" -o app ./cmd/main

FROM gcr.io/distroless/base-debian12
WORKDIR /app
COPY --from=builder /app/app .

EXPOSE 8080
CMD ["./app"]
