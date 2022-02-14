FROM alpine

RUN apk update  && apk add --no-cache ca-certificates

CMD ["/bin/sh", "-c", "./serverTCP"]

WORKDIR /tcpserver

COPY serverTCP .

RUN chmod +x serverTCP