# shorturl-go
[![Go Report Card](https://goreportcard.com/badge/github.com/korde96/shorturl-go)](https://goreportcard.com/report/github.com/korde96/shorturl-go)  
A basic URL Shortener Service  



## Build & Run
To build the project
```bash
go build ./cmd/shorturl-serve
```
This should generate the executable, which you can just run
```bash
./shorturl-serve
```  
The service requires a aerospike server running  
The config for which can be provided in ```config.json```
