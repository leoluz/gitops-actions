FROM golang:1.18-alpine3.15 AS builder

RUN mkdir /app
WORKDIR /app

ADD go.mod ./
ADD go.sum ./
RUN go mod download

ADD cmd ./cmd
RUN CGO_ENABLED=1 GOOS=linux go build ./cmd/goa
RUN ls -la


FROM alpine:3.15

RUN apk update && apk add git

RUN mkdir /app
WORKDIR /app

COPY --from=builder /app/goa /usr/local/bin/

CMD ["/usr/local/bin/goa"]
