FROM alpine:latest

WORKDIR /app
RUN apk add gcompat

COPY lookout .

ENTRYPOINT ["/app/lookout"]
