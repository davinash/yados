FROM ubuntu:latest
LABEL maintainer="dongre.avinash@gmail.com"

RUN apt-get update

COPY out/yadosctl /usr/bin

VOLUME ["/data"]
