FROM alpine:3.14

RUN apk update && apk upgrade
RUN apk add bind-tools curl bash bash-doc bash-completion

ENTRYPOINT ["tail", "-f", "/dev/null"]