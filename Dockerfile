FROM ubuntu:24.10

RUN apt-get update
RUN apt-get install -y texlive-full 
RUN apt-get install -y dvipng
RUN apt-get install -y golang-go
