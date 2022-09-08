# Start from golang base image
FROM golang:1.18 as builder

RUN apt-get update -y
WORKDIR /app

RUN wget https://github.com/wkhtmltopdf/packaging/releases/download/0.12.6.1-2/wkhtmltox_0.12.6.1-2.bullseye_amd64.deb -O wkhtmltopdf.deb

RUN apt-get install -f ./wkhtmltopdf.deb -y

COPY . .
