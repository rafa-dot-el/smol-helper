# syntax=docker/dockerfile:1.4
FROM alpine:3

RUN apk --no-cache --no-progress add git ca-certificates tzdata \
    && rm -rf /var/cache/apk/*

COPY smol-helper /

ENTRYPOINT ["/smol-helper"]
EXPOSE 80
