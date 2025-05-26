FROM ubuntu:noble

ARG DEBIAN_FRONTEND=noninteractive

RUN apt-get update && \
    apt-get install -y curl wget gpg ca-certificates golang-go libglib2.0-0t64\
    libnss3\
    libnspr4\
    libdbus-1-3\
    libatk1.0-0t64\
    libatk-bridge2.0-0t64\
    libatspi2.0-0t64\
    libxcomposite1\
    libxdamage1\
    libxext6\
    libxfixes3\
    libxrandr2\
    libgbm1\
    libxkbcommon0\
    libasound2t64          

WORKDIR /app

COPY go.mod ./
RUN go mod download
RUN go install github.com/playwright-community/playwright-go/cmd/playwright@latest
RUN /root/go/bin/playwright install

COPY . .

RUN go build -o scraper main.go

ENTRYPOINT ["/bin/bash", "entrypoint.sh"]

