FROM golang:1.24.3-alpine3.21 AS builder
WORKDIR /app
COPY . .
RUN go build -o main /app/main.go

FROM alpine:3.21
WORKDIR /app/
COPY --from=builder /app/main /app/
COPY app.env /app/
COPY db/migration ./db/migration
COPY entry-point.sh .
EXPOSE 4000
CMD [ "/app/main" ]
ENTRYPOINT ["/app/entry-point.sh"]
