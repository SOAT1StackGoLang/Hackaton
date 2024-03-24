# backend

## Description

This directory contains the Golang code for the backend of our application.

## Local Development with Docker-Compose

To run the backend locally, you can use the following command:

```bash
docker-compose build
docker-compose up -d
```

Client(Experimental) load testing can be done with the following command:

```bash
go build -o loadtest cmd/client/main.go
# edit the .testenv file with the desired number of requests
./loadtest
```
