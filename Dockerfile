FROM golang:latest

MAINTAINER Chris Purta <cpurta@gmail.com>

RUN apt-get update && \
    mkdir -p /app

ADD main.go boostrap.sh /app

WORKDIR /app

RUN bash boostrap.sh
