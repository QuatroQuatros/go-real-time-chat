FROM golang:1.24.5 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /chat-server ./cmd/server
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /seeder ./cmd/seeder

FROM gcr.io/distroless/base

WORKDIR /

COPY --from=builder /chat-server /chat-server
COPY --from=builder /seeder /seeder
COPY .env .

EXPOSE 8080

USER nonroot:nonroot

CMD ["./chat-server"]
