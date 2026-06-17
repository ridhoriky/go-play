FROM golang:1.26-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o go-app ./src/cmd/server/main.go

FROM alpine:3.22.4
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/go-app ./app
COPY --from=builder /app/config.yaml .
CMD ["./app"]