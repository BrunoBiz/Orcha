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

-- TBD
- Intended for Linux, tested on Debian 13

## Configuration

The full list of environment variables, descriptions, and examples can be found here:
- [configuration-reference](docs/configuration-reference.md)

## Installation

-- TBD

## Usage

GET  /containers\
GET  /containers/{id}\
GET  /containers/{id}/status\



POST /containers/server/{id}/start\
POST /containers/server/{id}/stop\
POST /containers/server/{id}/restart\
POST /containers/server/{id}/details\


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
- Game server management currently requires the Linux user under which the game server is installed to be provided as a request parameter
- Containerization and AWS deployment are planned
- A web frontend is planned
- Authentication is not yet implemented

## License

MIT
