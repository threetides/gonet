# gonet

Scaffold a Go HTTP service in seconds.

`gonet` is a small CLI that generates a ready-to-run Go HTTP project so you can
skip the boilerplate and start writing handlers. One command gives you a
standard folder layout, an initialized Go module, an optional git repository,
and a tiny `httpx` package with a shared HTTP client and JSON helpers.

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

The server starts on port `8080` with a health endpoint at
`GET /api/health`.

## What gets generated

```
my-service/
├── main.go                  # HTTP server with a /api/health endpoint
├── internal/
│   └── httpx/
│       ├── client.go        # shared *http.Client (30s timeout)
│       ├── helpers.go       # WriteJson helper
│       └── types.go         # standard JSON Response type
├── makefile                 # `make dev` runs the linter and the server
├── .gitignore
├── .env
└── go.mod
```

### The `httpx` package

`internal/httpx` is a minimal toolkit for building JSON APIs:

- **`Client`** — a shared `*http.Client` with a 30s timeout for outbound
  requests.
- **`WriteJson(w, statusCode, message, data)`** — writes a JSON response using
  a consistent envelope.
- **`Response`** — the JSON envelope (`message` plus optional `data`).

Example handler:

```go
mux.HandleFunc("GET /api/hello", func(w http.ResponseWriter, r *http.Request) {
    httpx.WriteJson(w, http.StatusOK, "Hello", map[string]string{"name": "world"})
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
