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
    texlive-full \
    dvipng && \
    rm -rf /var/lib/apt/lists/*

RUN sed -i 's/openin_any =.*/openin_any = p/' /usr/share/texlive/texmf-dist/web2c/texmf.cnf

WORKDIR /app
COPY --from=builder /app/vktexbot /app/vktexbot
COPY configs/config.yaml /app/config.yaml
COPY configs/secret.yaml /app/secret.yaml

CMD ["/app/vktexbot", "-config", "/app/config.yaml", "-secret", "/app/secret.yaml"]
