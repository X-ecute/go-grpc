FROM golang:1.24.1 as builder
RUN mkdir /app
ADD . /app
WORKDIR /app
RUN CGO_ENABLED=0 GOOS=linux go build -o app cmd/server/main.go

FROM alpine:latest AS production
# Install PostgreSQL client and bash
RUN apk add --no-cache postgresql-client bash
# Copy the built binary
COPY --from=builder /app/app .
# Copy migrations directory
COPY --from=builder /app/migrations ./migrations
# Copy wait script
COPY --from=builder /app/wait-for-db.sh .
RUN chmod +x wait-for-db.sh
# Copy any other necessary files
COPY --from=builder /app/.env .
CMD ["./wait-for-db.sh", "db", "./app"]