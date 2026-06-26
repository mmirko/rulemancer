FROM ubuntu:latest

RUN  mkdir -p /config

COPY rulemancer /rulemancer

EXPOSE 3000

CMD ["/rulemancer"]
