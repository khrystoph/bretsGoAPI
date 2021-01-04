FROM golang:latest

LABEL maintainer="Bret Edwards <sliverman69@gmail.com>"

ENV DOMAIN="example.com"
ENV APIKEY=""
ENV TESTING=""

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o bretsAPI .
ENTRYPOINT ["./bretsAPI"]