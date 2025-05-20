FROM ubuntu:24.10 AS builder

RUN apt-get update && \
    apt-get install -y --no-install-recommends \
    golang-go \
    ca-certificates \
    git

WORKDIR /app
COPY go.mod go.sum ./

RUN go mod download

COPY . .
RUN go build -o vktexbot cmd/main.go

FROM ubuntu:24.10

RUN apt-get update && \
    apt-get install -y --no-install-recommends \
    texlive-latex-base \
    texlive-latex-extra \
    texlive-latex-recommended \
    texlive-fonts-recommended \
    texlive-fonts-extra \
    texlive-lang-cyrillic \
    texlive-science \
    texlive-pictures \
    latexmk \
    dvipng && \
    rm -rf /var/lib/apt/lists/*

RUN sed -i 's/openin_any =.*/openin_any = p/' /usr/share/texlive/texmf-dist/web2c/texmf.cnf

WORKDIR /app
COPY --from=builder /app/vktexbot /app/vktexbot

CMD ["/app/vktexbot"]
