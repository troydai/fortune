FROM golang:1.18.2-alpine3.14 as builder

WORKDIR /usr/src/app

COPY go.mod ./
RUN go mod download && go mod tidy

COPY . .
RUN CGO_ENALBE=0 GOOS=linux go build -v -o /usr/local/bin/echoserver cmd/main.go

FROM alpine:3.14

RUN mkdir /app
COPY --from=builder /usr/local/bin/echoserver /app/echoserver

ENTRYPOINT ["/app/echoserver"]

