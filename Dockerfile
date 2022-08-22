FROM golang:1.19-alpine3.16 AS builder

RUN mkdir /app
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

ADD *.go ./
RUN go build -o hnrss .

FROM alpine:3.16

ENV GIN_MODE=debug

RUN apk add --no-cache ca-certificates

COPY --from=builder /app/hnrss /usr/local/bin

ENTRYPOINT [ "hnrss" ]

CMD [ "--bind=0.0.0.0:9000" ]

EXPOSE 9000
