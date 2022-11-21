FROM golang:1.19-alpine AS spire-agent-builder

RUN apk add build-base git mercurial

WORKDIR /src
RUN git clone --depth 1 --branch v1.5.1 https://github.com/spiffe/spire.git 

WORKDIR /src/spire
RUN go mod download
RUN make build

FROM golang:1.19-alpine AS builder

RUN apk add build-base git make

WORKDIR /src

COPY go.mod ./
RUN go mod download && go mod tidy

COPY . .
RUN make build

FROM alpine AS portal

RUN apk --no-cache add dumb-init ca-certificates curl
RUN mkdir -p /opt/spire/bin

COPY --from=builder /src/artifacts/portal /app/portal
COPY --from=spire-agent-builder /src/spire/bin/spire-agent /opt/spire/bin/spire-agent

ENTRYPOINT ["sleep", "infinity"]

FROM alpine AS front

RUN mkdir /app
COPY --from=builder /src/artifacts/front /app/front
COPY --from=spire-agent-builder /src/spire/bin/spire-agent /opt/spire/bin/spire-agent

USER 1000

ENTRYPOINT ["/app/front"]

FROM alpine AS datastore

RUN apk add --no-cache fortune
RUN mkdir /app
COPY --from=builder /src/artifacts/datastore /app/datastore
COPY --from=spire-agent-builder /src/spire/bin/spire-agent /opt/spire/bin/spire-agent

USER 1000

ENTRYPOINT ["/app/datastore"]
