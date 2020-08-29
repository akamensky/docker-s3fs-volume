FROM golang:1.15.0-buster AS builder

ADD . /build
RUN cd /build && go build .

FROM debian:buster-slim

RUN apt update && apt install s3backer fuse -y
COPY --from=builder /build/docker-s3fs-volume /bin/s3driver

CMD ["/bin/s3driver"]
