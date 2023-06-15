# syntax=docker/dockerfile:1
FROM golang:latest
WORKDIR /go/src/app

EXPOSE 8080
COPY go.mod go.sum Makefile env.sh ./
RUN go mod download
ARG MONGODB_URI

COPY cmd/ cmd/
COPY internal/ internal/

CMD ["make", "run"]
