FROM golang:1.15.0 AS builder
WORKDIR /go/src/github.com/akurczyk/golang_echo_blogging_platform/001_memory/app/src/
RUN go get -d -v github.com/labstack/echo && \
    go get -d -v github.com/go-playground/validator && \
    go get -d -v github.com/dgrijalva/jwt-go
COPY src .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

FROM alpine:3.12.0
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /go/src/github.com/akurczyk/golang_echo_blogging_platform/001_memory/app/src/main .
RUN adduser --disabled-password unprivileged && \
    chown -R root:unprivileged /app && \
    chmod -R 750 /app
USER unprivileged
EXPOSE 1323
CMD ["./main"]
