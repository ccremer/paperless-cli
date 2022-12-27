# paperless-cli
CLI tool to interact with paperless-ngx remote API


## Development

### Requirements

- go
- docker
- docker-compose (if running local test instance of paperless-ngx)
- goreleaser (if building deb/rpm packages locally)

### Build

Run `go run . --help` to directly invoke the CLI for testing purposes.
Run `make help` to see a list of available targets.

Commonly used:

- `make build`: Build the project
- `make local-install`: Start paperless-ngx in docker-compose (`http://localhost:8008`, user `admin:admin`)
