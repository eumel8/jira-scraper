FROM ghcr.io/mcsps/golang:1.0.8

WORKDIR /app

COPY go.mod ./
RUN go mod download
RUN go install github.com/playwright-community/playwright-go/cmd/playwright@latest
RUN /go/bin/playwright install

COPY . .

RUN go build -o scraper main.go

ENTRYPOINT ["/bin/bash", "entrypoint.sh"]

