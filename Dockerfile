FROM golang:latest AS build

ENV GO111MODULE on
ENV CGO_ENABLED 0

WORKDIR /go/src/github.com/ieee0824/valveMQ

COPY . .

RUN set -e \
    && go build -o /tmp/vmq cmd/vmq/main.go

FROM alpine:latest

RUN adduser -S vmq

COPY --from=build /tmp/vmq /bin/vmq

USER vmq

WORKDIR /home/vmq

CMD ["vmq"]