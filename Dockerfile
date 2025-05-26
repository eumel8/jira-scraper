FROM ubuntu:noble

ARG DEBIAN_FRONTEND=noninteractive

RUN apt-get update && \
    apt-get install -y curl wget gpg ca-certificates golang-go

WORKDIR /app

COPY go.mod ./
RUN go mod download
RUN go install github.com/playwright-community/playwright-go/cmd/playwright@latest
RUN /root/go/bin/playwright install

COPY . .

RUN go build -o scraper main.go

ENTRYPOINT ["/bin/bash", "entrypoint.sh"]

