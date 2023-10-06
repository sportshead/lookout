FROM golang:1.21 AS builder

WORKDIR /build

COPY go.* .
RUN go mod download

COPY *.go .
RUN go build -o lookout .

FROM alpine:latest AS app

WORKDIR /app

COPY --from=builder /build/lookout .

CMD ["./lookout"]
