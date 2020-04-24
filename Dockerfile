FROM golang:1.14-alpine3.11 AS build

RUN apk add --no-cache git

ADD . /go/action-tag
WORKDIR /go/action-tag

RUN go mod download
RUN go build -o action-tag github.com/akabos/action-tag

FROM alpine:3.11

RUN apk add --no-cache dumb-init
COPY --from=build /go/action-tag/action-tag /usr/local/bin

ENTRYPOINT ["/usr/bin/dumb-init", "--"]
CMD ["action-tag"]
