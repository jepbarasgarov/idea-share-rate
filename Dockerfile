# syntax=docker/dockerfile:1
FROM golang:1.17-buster as build
WORKDIR /go/src/app
ENV GOOS=linux
ENV GOARCH=amd64
COPY go.mod go.sum ./
RUN go mod download
RUN go mod verify
COPY . .
RUN go build daemon.go

FROM alpine:3.14 as prod
WORKDIR /home/belli
RUN apk update && apk add --no-cache nginx supervisor libc6-compat && rm /etc/nginx/http.d/default.conf
COPY --from=build /go/src/app/daemon ./
COPY ./deploy/nginx/app.conf /etc/nginx/http.d/
COPY ./deploy/supervisord.conf /etc/supervisord.conf
CMD ["/usr/bin/supervisord", "-c", "/etc/supervisord.conf"]