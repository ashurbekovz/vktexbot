FROM ubuntu:24.10

RUN apt-get update
RUN apt-get install -y texlive-full 
RUN apt-get install -y dvipng
RUN apt-get install -y golang-go

RUN sed -i 's/openin_any =*/openin_any = p/' /usr/share/texlive/texmf-dist/web2c/texmf.cnf && mktexlsr
RUN mktextfm larm0700
RUN mktextfm larm1000

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
WORKDIR /
