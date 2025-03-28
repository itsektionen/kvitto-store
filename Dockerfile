FROM golang:1.24.1 AS builder
COPY . /app
RUN cd /app && go mod verify && go build

FROM debian:bookworm
COPY --from=builder /app/kvitto-store /usr/bin/kvitto-store

CMD ["/usr/bin/kvitto-store"]