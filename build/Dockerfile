# syntax=docker/dockerfile:1
FROM golang:1.20-alpine
WORKDIR /go/src/app

EXPOSE 8080
COPY go.mod go.sum ./
RUN go mod download

COPY cmd/ cmd/
COPY internal/ internal/
RUN go build .

CMD ["./solvify-backend"]
