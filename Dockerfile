FROM golang:1.24.3-alpine3.21 AS builder
WORKDIR /app
COPY . .
RUN go build -o main /app/main.go
RUN apk add curl
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.18.3/migrate.linux-amd64.tar.gz | tar xvz

FROM alpine:3.21
WORKDIR /app/
COPY --from=builder /app/main /app/
COPY --from=builder /app/migrate .
COPY app.env /app/
COPY db/migration ./migration
COPY entry-point.sh .
EXPOSE 4000
CMD [ "/app/main" ]
ENTRYPOINT ["/app/entry-point.sh"]
