FROM golang:1.20 AS builder
WORKDIR /app
ADD . /app
RUN [ "make" ]

FROM alpine
COPY --from=builder /app/redis-ha-check /app/redis-ha-check
ENTRYPOINT [ "/app/redis-ha-check" ]
