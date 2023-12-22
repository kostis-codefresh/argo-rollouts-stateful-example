FROM golang:1.21-alpine3.18 AS build-env

WORKDIR /tmp/workdir

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build 

FROM alpine:3.18

EXPOSE 8080

RUN apk add --no-cache ca-certificates bash

COPY --from=build-env /tmp/workdir/static /app/static
COPY --from=build-env /tmp/workdir/interest /app/interest

WORKDIR /app

CMD ["./interest"]
