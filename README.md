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

curl -X POST http://localhost:8080/api/inboxes \
  -H 'Content-Type: application/json' \
  -d '{"name":"Support"}'

curl http://localhost:8080/api/inboxes
```

HookLens creates its SQLite database at the operating system's user config
directory, under `HookLens/hooklens.db`. The location is resolved at runtime
with Go's `os.UserConfigDir()`, so the same binary source works across macOS,
Linux, and Windows.

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
