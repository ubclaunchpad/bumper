# bumper
bump other players and objects off to stay alive

[![Go Report Card](https://goreportcard.com/badge/github.com/ubclaunchpad/bumper)](https://goreportcard.com/report/github.com/ubclaunchpad/bumper) [![Deployed with Inertia](https://img.shields.io/badge/deploying%20with-inertia-blue.svg)](https://github.com/ubclaunchpad/inertia)

![gif](/.static/bumper.gif)

## Getting Started
Play the latest version [here](http://bumper.ubclaunchpad.com)!  

## Quickstart

Go and Node are required. To install the required dependencies:

```bash
$ make deps
```

Create a .env file with the following variables and place it in the project root:
```
NODE_ENV=$YOUR_VAR
SERVER_URL=$YOUR_VAR
DATABASE_URL=$YOUR_VAR
PORT=$YOUR_VAR
```

### Run the Server

```bash
$ make server
```

### Start the Client

```bash
$ make client
```

Play at localhost:8080!
