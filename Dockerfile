FROM golang:latest AS builder
ENV GOPROXY=https://goproxy.cn
ENV GO111MODULE=on
WORKDIR /spblog
COPY . /spblog/
RUN go mod download \
&& CGO_ENABLED=0 GOOS=linux go build -o spblog .

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /spblog/spblog .
COPY ./conf/conf.yaml ./conf/conf.yaml
COPY ./static ./static
COPY ./views ./views
RUN ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
EXPOSE 8090
CMD ["./spblog"]