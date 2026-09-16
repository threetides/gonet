# gonet

A CLI that scaffolds a Go HTTP service.

`gonet` generates a Go HTTP project: a standard folder layout, an initialized
Go module, an optional git repository, an `apperr` package for typed
application errors, and an `httpx` package with a shared HTTP client plus JSON
and error response helpers.

## Installation

### With `go install`

```sh
go install github.com/threetides/gonet@latest
```

### From a release

Download a prebuilt binary for your platform from the
[releases page](https://github.com/threetides/gonet/releases), extract it, and
put the `gonet` binary somewhere on your `PATH`.

## Usage

Create a new project in a new directory:

```sh
gonet init my-service
```

Scaffold into the current directory:

```sh
gonet init .
```

Or run it with no arguments and you'll be prompted for a name:

```sh
gonet init
```

During init you'll be asked for a module name (defaulting to the project name)
and whether to initialize a git repository.

Once it's done:

```sh
cd my-service
go run main.go   # or: make dev
```

The server reads `PORT` from `.env` (default `8080`) and listens on
`127.0.0.1`. It exposes a health endpoint at `GET /api/health`.

## What gets generated

```
my-service/
├── main.go                  # HTTP server with a /api/health endpoint
├── internal/
│   ├── apperr/
│   │   └── apperr.go        # Apperr type, Kind enum and constructor helpers
│   └── httpx/
│       ├── client.go        # shared *http.Client (30s timeout)
│       ├── helpers.go       # WriteJSON and WriteError helpers
│       └── types.go         # standard JSON Response type
├── makefile                 # `make dev` runs the server
├── .gitignore
├── .env                     # PORT=8080
└── go.mod
```

### The `apperr` package

`internal/apperr` defines a single error type, `Apperr`, that carries a
`Kind`, a user-facing `Message`, an optional wrapped `Err` and optional
`Data`. Each `Kind` maps to an HTTP status via `Kind.HTTPStatus()`.
Constructor helpers cover the common cases:

- `InternalServerError(err)` wraps an unexpected error. It is also the zero
  value, so any `Kind` you forget to set is treated as a 500.
- `BadRequest(data)`, `Conflict(data)` attach details for the client.
- `Unauthorized()`, `Forbidden()`, `NotFound()`, `TooManyRequests()`,
  `ServiceUnavailable()` return a plain message.

### The `httpx` package

`internal/httpx` holds a few helpers for building JSON APIs:

- **`Client`**: a shared `*http.Client` with a 30s timeout for outbound
  requests.
- **`WriteJSON(w, statusCode, message, data)`**: writes a JSON response using a
  consistent envelope.
- **`WriteError(w, r, err)`**: unwraps an `*apperr.Apperr` and writes its
  status, message and data. Anything that is not an `Apperr`, or is an
  internal server error, is logged with the request path and returned as a
  generic 500 so internals never leak to the client.
- **`Response`**: the JSON envelope (`message` plus optional `data`).

Example handler, where the service returns an `apperr` such as
`apperr.NotFound()`:

```go
mux.HandleFunc("GET /api/users/{id}", func(w http.ResponseWriter, r *http.Request) {
    user, err := users.Get(r.PathValue("id"))
    if err != nil {
        httpx.WriteError(w, r, err)
        return
    }
    httpx.WriteJSON(w, http.StatusOK, "success", user)
})
```

## Development

`gonet` is built with [Cobra](https://github.com/spf13/cobra). The generated
files live in [`internal/templates/`](internal/templates/) and are embedded
into the binary at build time.

```sh
go build ./...   # build
go run . init    # try it out
```

Releases are cut with [GoReleaser](https://goreleaser.com) via the GitHub
Actions workflow in [`.github/workflows/`](.github/workflows/).

## License

See [LICENSE](LICENSE).
