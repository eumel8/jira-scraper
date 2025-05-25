FROM mcr.microsoft.com/playwright:v1.44.0-focal

WORKDIR /app

COPY go.mod ./
RUN apt-get update && apt-get install -y golang-go
RUN go mod download

COPY . .

RUN go build -o scraper main.go

ENTRYPOINT ["/bin/bash", "entrypoint.sh"]

