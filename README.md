# paperless-cli

CLI tool to interact with paperless-ngx remote API

## Subcommands

- `upload`: Uploads local document(s) to Paperless instance
- `consume`: Consumes a local directory and uploads each file to Paperless instance. The files will be deleted once uploaded.
- `bulk-download`: 

## Installation

Go:
`go install github.com/ccremer/paperless-cli@latest`

Docker:
`docker run ghcr.io/ccremer/paperless-cli:latest`

Binary:
```bash
wget https://github.com/ccremer/paperless-cli/releases/latest/download/paperless-cli_linux_amd64
chmod +x paperless-cli_linux_amd64
sudo mv paperless-cli_linux_amd64 /usr/local/bin/paperless-cli
```

Deb:
```bash
wget https://github.com/ccremer/paperless-cli/releases/latest/download/paperless-cli_linux_amd64.deb
sudo dpkg -i paperless-cli_linux_amd64.deb
rm paperless-cli_linux_amd64.deb
```

RPM:
```bash
wget https://github.com/ccremer/paperless-cli/releases/latest/download/paperless-cli_linux_amd64.rpm
sudo rpm -i paperless-cli_linux_amd64.rpm
rm paperless-cli_linux_amd64.rpm
```

## Systemd Service

The `consume` subcommand is a long-running process that is best run as a daemon.
The Deb/RPM packages come with a SystemD unit file.

Enable SystemD `consume` service:
```bash
sudo ${EDITOR:-nano} /etc/default/paperless-cli
sudo systemctl enable paperless-consume
sudo systemctl start paperless-consume
```

## Why does this exist?

I didn't find any other projects or means to consume a directory that _uploads_ the documents via API.
In my case, I can't configure the scanner to directly upload to the consume dir as setup by paperless-already, I have to watch the dir on a different host.
So I created a tool that also watches a directory, but uploads them to Paperless instead.

Other projects that I've found:

- https://github.com/stgarf/paperless-cli (archived, doesn't upload or consume)

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
