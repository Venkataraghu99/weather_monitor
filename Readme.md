## Weather Monitor

This service is implemented using gorilla mux router and environment variables are defined in the .env file

```bash
go run cmd/main.go
```

```bash
curl --location 'http://localhost:8080/weather/status?latitude=44.34&longitude=10.99'
```