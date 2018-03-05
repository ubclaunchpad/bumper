# bumper
bump other players and objects off to stay alive

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
