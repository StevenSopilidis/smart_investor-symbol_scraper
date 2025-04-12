FROM golang:1.23-alpine AS Build

WORKDIR /app
COPY . .
RUN go mod download

RUN go build -o /symbol-scraper main.go

FROM alpine:latest

WORKDIR /
COPY --from=Build /symbol-scraper /symbol-scraper

ENTRYPOINT [ "/symbol-scraper" ]