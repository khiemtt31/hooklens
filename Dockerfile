# syntax=docker/dockerfile:1

FROM golang:1.26.5-bookworm AS build

WORKDIR /src

COPY go.mod go.sum ./
COPY cmd ./cmd
COPY internal ./internal

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -trimpath -ldflags="-s -w" \
    -o /out/hooklens ./cmd/hooklens

FROM gcr.io/distroless/static-debian12:nonroot

COPY --from=build /out/hooklens /hooklens

EXPOSE 8080

ENTRYPOINT ["/hooklens"]
