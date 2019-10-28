FROM golang:latest

ENV GOPROXY=https://goproxy.cn
ENV GO111MODULE=on

MAINTAINER Razil "shenpu1819@163.com"

WORKDIR /spblog

COPY . /spblog/
RUN go mod download
RUN go build -o spblog main.go
RUN chmod 777 /spblog/spblog

EXPOSE 80

ENTRYPOINT ["./spblog"]
