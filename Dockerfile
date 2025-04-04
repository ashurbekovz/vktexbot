FROM ubuntu:24.10

RUN apt-get update
RUN apt-get install -y texlive-full 
RUN apt-get install -y dvipng
RUN apt-get install -y golang-go
RUN apt-get install -y bubblewrap
RUN sed -i 's/openin_any =*/openin_any = p/' /usr/share/texlive/texmf-dist/web2c/texmf.cnf && mktexlsr
RUN apt-get clean && \
    rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
WORKDIR /
