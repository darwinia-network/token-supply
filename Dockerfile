FROM ubuntu:18.04

WORKDIR /app
COPY ./config ./config
COPY ./token ./token

RUN apt update -y && apt install curl -y

EXPOSE 5344
CMD ["/app/token"]
