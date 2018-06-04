# bumper
bump other players and objects off to stay alive

[![Go Report Card](https://goreportcard.com/badge/github.com/ubclaunchpad/bumper)](https://goreportcard.com/report/github.com/ubclaunchpad/bumper) [![Deployed with Inertia](https://img.shields.io/badge/Deploying%20with-Inertia-blue.svg)](https://github.com/ubclaunchpad/inertia)

## Getting Started
Play the latest version [here](http://bumper.ubclaunchpad.com)!  

## Docker Quickstart
Install Docker and the Docker Compose toolset. Then run:

```bash
make docker-start
```

## Manual Quickstart

Go and Node are required. To install the required dependencies:

```bash
$ make deps
```

### Run the Server

```bash
$ make server
```

To add dependencies:

```bash
$ dep ensure -add github.com/my/dependency
```

### Start the Client

```bash
$ make client
```

Play at localhost:8080!
