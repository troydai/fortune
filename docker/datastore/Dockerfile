FROM golang:1.18.2-alpine3.14 as builder

WORKDIR /usr/src/app

COPY go.mod ./
RUN go mod download && go mod tidy

COPY . .
RUN CGO_ENALBE=0 GOOS=linux go build -v -o /usr/local/bin/datastore cmd/datastore/main.go

FROM alpine:3.14

RUN apk add --no-cache fortune
RUN mkdir /app
COPY --from=builder /usr/local/bin/datastore /app/datastore

USER 1000

ENTRYPOINT ["/app/datastore"]
