FROM golang:1.18 

WORKDIR /usr/src/app

COPY go.mod ./
RUN go mod download && go mod tidy

COPY . .
RUN go build -v -o /usr/local/bin/echoserver cmd/main.go

ENTRYPOINT ["echoserver"]

