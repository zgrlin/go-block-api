## Go Ethereum Block API

Golang implementation of Ethereum Block API

This API can store (in memory) new block information by getting Etherscan.

## Running

```shell
$ go mod init
$ go run main.go
```

## Manual POST

```shell
$ curl -X POST http://127.0.0.1:8080/blocks -H "Content-Type: application/json" -d '{"height": "", "date": "", "mined": "", "reward": "", "difficult": "}'  
```

## GET

```shell
$ curl -X GET http://127.0.0.1:8080/blocks
```
