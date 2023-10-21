FROM golang:1.21-alpine as builder

WORKDIR /identifier_service

COPY ../.. ./


RUN go mod tidy && go mod download

RUN go build -o app ./cmd/main.go


FROM alpine

WORKDIR /

COPY --from=builder /identifier_service/app /app
COPY --from=builder /identifier_service/Deployment/app/app.env /app.env
COPY --from=builder /identifier_service/internal/repository/postgres/migrations/ /migrations

CMD ["/app"]