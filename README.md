# bumper
bump other players and objects off to stay alive

[![Go Report Card](https://goreportcard.com/badge/github.com/ubclaunchpad/bumper)](https://goreportcard.com/report/github.com/ubclaunchpad/bumper) [![Deployed with Inertia](https://img.shields.io/badge/Deploying%20with-Inertia-blue.svg)](https://github.com/ubclaunchpad/inertia)

# Getting Started
How to get started running a server and client on your local machine.

## docker-compose Startup

Install Docker and the docker-compose toolsets. Then run:

```bash
$ docker-compose -f docker-compose.dev.yml up
```

This will start up the server and client.

## Manual Startup
### Running the server
1. Enter the `bumper/server/` directory
2. Get the websocket library `go get github.com/gorilla/websocket`
3. Run the server `go run main.go`

#### Server dependencies
[dep](https://github.com/golang/dep) is used for handling server dependencies.
```bash
$ cd ./server
$ go get -u github.com/golang/dep/cmd/dep
$ dep ensure
$ go run main.go
```
To add dependencies:
```bash
$ dep ensure -add github.com/my/dependency
```

### Running the client
1. Enter the `bumper/client/` directory
2. Install dependency packages `npm install`
3. Start and build the client `npm start`
4. Open a browser and connect to `localhost:8080`
5. Play!

# Development

## Client

```bash
$ npm install
```

## Server

```bash
$ go get -u github.com/golang/dep/cmd/dep
$ cd server ; dep ensure
```
