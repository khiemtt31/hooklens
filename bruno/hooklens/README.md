# HookLens Bruno collection

Open the `bruno/hooklens` directory in Bruno and select the `local` environment.

Start the API from the repository root first:

```bash
go run ./cmd/hooklens
```

Requests included:

- `Get health` — verifies the `200` health report.
- `Reject non-GET health requests` — verifies the `405` response for `POST`.
- `POST ping` — sends a `pong` string and verifies it is returned.
- `POST inbox` — creates a named inbox.
- `GET inboxes` — lists the stored inboxes.
