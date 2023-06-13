# syntax=docker/dockerfile:1
FROM golang:latest
WORKDIR /go/src/app

EXPOSE 8080
COPY go.mod go.sum Makefile ./
RUN go mod download

COPY cmd/ cmd/
COPY internal/ internal/

CMD ["make", "run"]
