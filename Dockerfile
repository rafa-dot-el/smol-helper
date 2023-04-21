FROM golang:1-alpine as builder

RUN apk --no-cache --no-progress add make git

WORKDIR /go/smol-helper

ENV GO111MODULE on

# Download go modules
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN go build -o dist/smol-helper cmd/smol-helper/main.go

FROM alpine:3
RUN apk update \
    && apk add --no-cache ca-certificates tzdata \
    && update-ca-certificates

COPY --from=builder /go/smol-helper/dist/smol-helper /usr/bin/smol-helper

ENTRYPOINT [ "/usr/bin/smol-helper" ]
