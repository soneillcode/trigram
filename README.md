
## Trigram Service

A service for storing text data and using to generate text using [Trigrams](https://en.wikipedia.org/wiki/Trigram).

## getting started

```
go mod tidy
```

run the service
```
go run cmd/main.go
```

run the service using docker
```
docker build -t dev_service ./service
docker run --rm -d --name dev_service -p 8080:8080 dev_service
```

enable debug
```
go build -gcflags='all=-N -l' cmd/main.go
```
