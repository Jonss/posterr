FROM golang:1.18 AS gobuilder
WORKDIR /app
COPY pkg/ pkg/
COPY api/ api/
COPY config/ config/
COPY db/ db/
COPY cmd/posterr/main.go /app/main.go
COPY .env /app/.env
COPY go.mod /app/go.mod
COPY go.sum /app/go.sum
RUN CGO_ENABLED=0 GOOS=linux go build -o bin main.go

FROM alpine:3.15.1 as posterr-app
COPY --from=gobuilder /app/bin bin
COPY --from=gobuilder /app/db/migrations migrations
COPY --from=gobuilder /app/.env .env
CMD ["bin"]