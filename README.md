# Orcha

**Orcha is a Go-based REST API for orchestrating a Proxmox environment and managing game servers running inside its LXC containers.**

It provides access to Proxmox infrastructure information while providing a standardized interface for managing game-server lifecycle operations independently of the underlying game-server manager.

## Overview

Orcha acts as an orchestration layer between clients, Proxmox, and game servers running inside LXC containers.

It communicates directly with the Proxmox API when retrieving infrastructure information and uses SSH through the Proxmox node to interact with services running inside containers.

The API is designed to remain agnostic to how an individual game server is managed. Server-specific management is delegated to the appropriate game-server manager, such as [Narwhal](https://github.com/BrunoBiz/Narwhal) or LinuxGSM.

## Features

* Retrieve Proxmox LXC container information and status
* Manage game-server lifecycle operations
* Support different game-server management implementations
* SSH-based interaction with LXC containers
* Structured logging with custom log levels
* Console and file logging
* Environment-based configuration
* Code-first OpenAPI documentation
* C4 architecture documentation
* Automated build and SSH-based deployment

## Architecture

Orcha is composed mainly of:
- A **Proxmox client** that wraps the PVE API to retrieve infrastructure information about the Proxmox environment
- An **SSH client** used to communicate with containers and manage game servers through the main Proxmox node
- A **web framework and router** responsible for handling API requests
- A **logging system** responsible for console and file logging, including custom log levels

The API communicates with Proxmox for infrastructure information and uses SSH to access services running inside the LXC containers.

For game-server operations, the API delegates management to the configured game-server manager rather than implementing game-specific behavior itself.

### C4 Diagrams
- [C4 Context Diagram](docs/architecture/c4-level1-system-context.png)
- [C4 Container Diagram](docs/architecture/c4-level2-container.png)

## Technologies

* **Go**
* **Gin**
* **Huma**
* **Proxmox API**
* **SSH**
* **Viper**
* **slog**
* **OpenAPI**
* **C4 Model / PlantUML**

## Requirements

### Orcha Host

Orcha is intended to run on Linux and is currently developed and tested on
Debian 13.

To build and run Orcha from source, the host requires:

- Go 1.26 or newer
- Git
- Network access to the Proxmox VE API
- SSH access to the Proxmox node
- A Proxmox VE API token
- An SSH private key for connecting to the Proxmox node
- A `pm.env` configuration file

Orcha currently connects to a single configured Proxmox node.

The configuration file must be named:

```text
pm.env
```

and should be located in the directory from which Orcha is started.

See
[`docs/configuration-reference.md`](docs/configuration-reference.md)
for the complete list of configuration variables.

### Game Server Requirements

For game-server management, each managed LXC container must have a dedicated
Linux user named:

```text
gameserver
```

Orcha expects the game-server manager entry point to be available as:

```text
/home/gameserver/Narwhal
```

The manager must support the following operations:

```text
./Narwhal start
./Narwhal stop
./Narwhal restart
./Narwhal details
```

[Narwhal](https://github.com/BrunoBiz/Narwhal) implements this interface
directly.

LinuxGSM can also be used through a compatibility wrapper named `Narwhal`
placed in `/home/gameserver`. The wrapper forwards the lifecycle operation to
the appropriate LinuxGSM server script.

For example:

```bash
#!/bin/bash

echo "Using LinuxGSM - Wrapper"

exec ./gmodserver "$@"
```

The LinuxGSM installation and its wrapper must run under the `gameserver` user.

## Configuration

The full list of environment variables, descriptions, and examples can be found here:
- [configuration-reference](docs/configuration-reference.md)

## Installation

### 1. Install Git and Go

Install Git using your distribution's package manager.

For Debian-based systems:

```bash
sudo apt update
sudo apt install -y git
```

Verify that Go 1.26 or newer is installed:

```bash
go version
```

### 2. Clone the Repository

```bash
git clone https://github.com/BrunoBiz/Orcha.git
cd Orcha
```

### 3. Download Dependencies

```bash
go mod download
```

### 4. Build Orcha

```bash
go build -o Orcha
```

This creates the Orcha executable in the project directory:

```text
./Orcha
```

### 5. Configure Orcha

Create a file named:

```text
pm.env
```

in the project directory.

Example:

```env
PVE_URL=https://192.168.1.10:8006/api2/json
PVE_USER=Orcha
PVE_REALM=pve
PVE_USER_REALM=Orcha@pve
PVE_TOKEN_ID=api-token
PVE_TOKEN=your-proxmox-api-token
PVE_NODE_NAME=main
SSH_KEY_FILE=/path/to/private/key
SSH_KEY_PASSPHRASE=your-key-passphrase
SSH_PVE_IP=192.168.1.10
SSH_PVE_PORT=22
```

See
[`docs/configuration-reference.md`](docs/configuration-reference.md)
for a complete description of each variable.

### 6. Start Orcha

Run Orcha from the directory containing `pm.env`:

```bash
./Orcha
```

By default, the API listens on:

```text
http://localhost:8090
```

Interactive OpenAPI documentation is available at:

```text
http://localhost:8090/docs
```

The OpenAPI specification is also available at:

```text
/openapi.json
/openapi-3.0.json
/openapi.yaml
/openapi-3.0.yaml
```

## Usage

Orcha exposes the following API operations:

```text
GET  /containers
GET  /containers/{id}
GET  /containers/{id}/status

POST /containers/server/{id}/start
POST /containers/server/{id}/stop
POST /containers/server/{id}/restart
POST /containers/server/{id}/details
```

`{id}` represents the Proxmox VMID of the target LXC container.

For complete request and response documentation, use the interactive OpenAPI
documentation available at:

```text
/docs
```

## API Documentation

Orcha uses Huma to generate its OpenAPI documentation directly from the API definitions.

The generated OpenAPI specification is exposed by the running API and can be accessed at:

`/openapi.json`

`/openapi-3.0.json`

`/openapi.yaml`

`/openapi-3.0.yaml`

Interactive API documentation is available through:

`/docs`


## Deployment

Orcha is currently deployed to my self-hosted Proxmox environment using an automated Go build and SSH-based deployment workflow.


## Current Scope

The current version targets Proxmox LXC environments and supports game-server management through the supported game-server managers.

## Known Limitations

- The SSH client is currently connecting to the root user
- Manual and laborious first-time installation
- LinuxGSM needs a small bash wrapper for the API to remain agnostic
- Containerization and AWS deployment are planned
- A web frontend is planned
- Authentication is not yet implemented
- SSH host-key verification is currently disabled; the Proxmox SSH connection uses `ssh.InsecureIgnoreHostKey()` and is intended for use within a trusted environment.

## License

MIT
