ARG GO_VERSION=1.18.3

FROM golang:${GO_VERSION} as builder
WORKDIR /go/src/github.com/darwinia-network/token

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o /token

FROM buildpack-deps:buster-scm

WORKDIR /app
COPY ./config ./config
COPY --from=builder /token ./token

EXPOSE 5344
CMD ["/app/token"]
