FROM ubuntu:24.10

RUN apt-get update
RUN apt-get install -y \
    texlive-latex-base \
    texlive-latex-extra \
    texlive-latex-recommended \
    texlive-fonts-recommended \
    texlive-fonts-extra \
    texlive-lang-cyrillic \
    texlive-science \
    texlive-pictures
RUN apt-get install -y latexmk
RUN apt-get install -y dvipng
RUN apt-get install -y golang-go

RUN sed -i 's/openin_any =*/openin_any = p/' /usr/share/texlive/texmf-dist/web2c/texmf.cnf && mktexlsr

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
WORKDIR /
