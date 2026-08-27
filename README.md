# HookLens

## Run the API

From the repository root:

```bash
go run ./cmd/hooklens
```

The API starts at `http://localhost:8080`. Stop it with `Ctrl+C`.

## Try the endpoints

In a second terminal:

```bash
curl http://localhost:8080/api/v1/health

curl -X POST http://localhost:8080/ping \
  -H 'Content-Type: application/json' \
  -d '{"pong":"hello"}'
```

## Run tests and checks

```bash
go test ./...
go test -race ./...
go vet ./...
```

## Run with live reload

If [Air](https://github.com/air-verse/air) is installed:

```bash
air
```
