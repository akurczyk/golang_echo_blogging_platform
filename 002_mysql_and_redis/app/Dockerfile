FROM golang:1.15.0 AS builder
WORKDIR /go/src/github.com/akurczyk/golang_echo_blogging_platform/002_mysql_and_redis/app/src/
RUN go get -v github.com/swaggo/swag/cmd/swag
COPY src .
RUN /go/bin/swag init
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

FROM alpine:3.12.0
RUN apk --no-cache add ca-certificates
RUN mkdir /app
WORKDIR /app
COPY --from=builder /go/src/github.com/akurczyk/golang_echo_blogging_platform/002_mysql_and_redis/app/src/main ./
COPY misc/entrypoint.sh misc/wait-for-it.sh ./
RUN adduser --disabled-password unprivileged && \
    chown -R root:unprivileged /app && \
    chmod -R 750 /app
USER unprivileged
EXPOSE 1323
ENTRYPOINT ["./entrypoint.sh"]
