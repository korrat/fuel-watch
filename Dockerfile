FROM golang:alpine AS builder

RUN apk add --no-cache ca-certificates git dep

RUN go get -d github.com/korrat/fuel-watch
WORKDIR /go/src/github.com/korrat/fuel-watch

RUN dep ensure
RUN go build

FROM alpine:latest

WORKDIR /root/

RUN apk add --no-cache ca-certificates

COPY --from=builder /go/src/github.com/korrat/fuel-watch/fuel-watch .

CMD [ "./fuel-watch" ]