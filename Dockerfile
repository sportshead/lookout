FROM alpine:latest

WORKDIR /app

COPY lookout .

CMD ["./lookout"]
