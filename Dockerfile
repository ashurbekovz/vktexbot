FROM ubuntu:24.10

RUN apt-get update
RUN apt-get install -y texlive-full 
RUN apt-get install -y dvipng
RUN apt-get install -y golang-go

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
