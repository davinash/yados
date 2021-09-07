FROM ubuntu:latest
LABEL maintainer="dongre.avinash@gmail.com"

COPY out/yadosctl /usr/bin
VOLUME ["/data"]

