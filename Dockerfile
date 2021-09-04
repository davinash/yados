FROM ubuntu:latest
LABEL maintainer="dongre.avinash@gmail.com"

RUN apt-get update
RUN apt-get install -y wget
RUN apt-get install -y git
RUN apt-get install -y gcc
RUN apt-get install -y curl
RUN apt-get install -y make

RUN wget -q https://golang.org/dl/go1.17.linux-amd64.tar.gz
RUN rm -rf /usr/local/go && tar -C /usr/local -xzf go1.17.linux-amd64.tar.gz

ENV PATH="/usr/local/go/bin:${PATH}"

RUN git clone https://github.com/davinash/yados
RUN cd yados && make

VOLUME ["/data"]
