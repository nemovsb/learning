FROM alpine

RUN apk update  && apk add --no-cache ca-certificates

CMD ["/bin/sh", "-c", "./learning_for_alpine"]

EXPOSE 8081

WORKDIR /learning

COPY . /learning/

RUN chmod +x learning_for_alpine



