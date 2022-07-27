FROM ubuntu:18.04

WORKDIR /app
COPY ./config ./config
COPY ./token ./token

EXPOSE 5344
CMD ["/app/token"]
