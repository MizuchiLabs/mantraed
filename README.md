<p align="center">
<img src="https://raw.githubusercontent.com/MizuchiLabs/mantrae/refs/heads/main/.github/logo.svg" width="80">
<br><br>
<img alt="GitHub Tag" src="https://img.shields.io/github/v/tag/MizuchiLabs/mantraed?label=Version">
<img alt="GitHub License" src="https://img.shields.io/github/license/MizuchiLabs/mantraed">
<img alt="GitHub Issues or Pull Requests" src="https://img.shields.io/github/issues/MizuchiLabs/mantraed">
</p>

# Mantraed

**Mantraed** (Mantræ daemon) is the agent component for [Mantræ](https://github.com/MizuchiLabs/mantrae), a web-based configuration manager for Traefik. Mantraed monitors Docker containers in real-time and automatically syncs their Traefik labels to your Mantræ server.

## Features

- **Real-time Monitoring**: Watches Docker containers for changes instantly
- **Automatic Sync**: Sends container labels to Mantræ server automatically
- **Multi-host Support**: Deploy agents on multiple machines to manage distributed containers
- **Lightweight**: Minimal resource footprint
- **Docker Socket Access**: Reads container information directly from Docker

## How It Works

Mantraed connects to the Docker socket on your host machine and monitors container events. When containers are created, updated, or removed, it extracts their Traefik labels and sends them to your Mantræ server via the configured API endpoint.

## Quick Start

### Prerequisites

- Docker running on your host machine
- A running [Mantræ](https://github.com/MizuchiLabs/mantrae) server
- An agent token from your Mantræ server

### Docker Compose (Recommended)

```yaml
services:
  mantraed:
    image: ghcr.io/mizuchilabs/mantraed:latest
    container_name: mantraed
    network_mode: host # Required for detecting the hostname
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock
    environment:
      - TOKEN=<your-agent-token>
      - HOST=https://mantrae.example.com
    restart: unless-stopped
```

### Docker Run

```bash
docker run -d \
  --name mantraed \
  --network host \
  -v /var/run/docker.sock:/var/run/docker.sock \
  -e TOKEN=<your-agent-token> \
  -e HOST=https://mantrae.example.com \
  --restart unless-stopped \
  ghcr.io/mizuchilabs/mantraed:latest
```

### Binary Installation

Download the latest release from [releases](https://github.com/MizuchiLabs/mantraed/releases) and run:

```bash
export TOKEN=<your-agent-token>
export HOST=https://mantrae.example.com
./mantraed
```

## Configuration

### Environment Variables

| Variable      | Required | Description                                                     |
| ------------- | -------- | --------------------------------------------------------------- |
| `TOKEN`       | Yes      | Authentication token from Mantræ server                         |
| `HOST`        | Yes      | URL of your Mantræ server (e.g., `https://mantrae.example.com`) |
| `DOCKER_HOST` | No       | Docker socket path (default: `unix:///var/run/docker.sock`)     |

### Getting an Agent Token

1. Log into your Mantræ server
2. Navigate to Settings → Agents
3. Create a new agent profile
4. Copy the generated token

## Usage Example

Once mantraed is running, simply add Traefik labels to your containers:

```yaml
services:
  myapp:
    image: myapp:latest
    labels:
      - traefik.enable=true
      - traefik.http.routers.myapp.rule=Host(`myapp.example.com`)
      - traefik.http.routers.myapp.entrypoints=websecure
      - traefik.http.services.myapp.loadbalancer.server.port=8080
```

Mantraed will automatically detect these labels and sync them to your Mantræ server.

## Multi-host Deployment

Deploy mantraed on each host where you run containers:

```yaml
# Host 1
services:
  mantraed:
    image: ghcr.io/mizuchilabs/mantraed:latest
    network_mode: host
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock
    environment:
      - TOKEN=<token-for-host1>
      - HOST=https://mantrae.example.com
    restart: unless-stopped
```

Each agent will report its hostname and containers independently to the Mantræ server.

## Troubleshooting

### Agent not connecting to Mantræ

- Verify the `HOST` URL is accessible from the agent
- Check that the `TOKEN` is valid and not expired
- Ensure network connectivity between agent and server

### Containers not being detected

- Verify Docker socket is mounted correctly (`/var/run/docker.sock`)
- Check that the agent has permissions to read the Docker socket
- Ensure `network_mode: host` is set for proper hostname detection

### Labels not syncing

- Verify containers have `traefik.enable=true` label
- Check agent logs for any errors
- Confirm the token has appropriate permissions on the Mantræ server

## Documentation

For more information about Mantræ and its ecosystem, visit the [main documentation](https://mizuchilabs.github.io/mantrae/).

## Contributing

Contributions are welcome! Feel free to submit issues, fork the repository, and create pull requests.

## License

MIT License - See [LICENSE](LICENSE)

## Related Projects

- [**Mantræ**](https://github.com/MizuchiLabs/mantrae) - The main configuration manager
- [**Traefik**](https://traefik.io/) - The reverse proxy this agent supports
