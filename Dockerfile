FROM alpine:latest

WORKDIR /app

COPY lookout .

ENTRYPOINT ["/app/lookout"]
