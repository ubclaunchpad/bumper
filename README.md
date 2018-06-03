# bumper
bump other players and objects off to stay alive

[![Go Report Card](https://goreportcard.com/badge/github.com/ubclaunchpad/bumper)](https://goreportcard.com/report/github.com/ubclaunchpad/bumper) [![Deployed with Inertia](https://img.shields.io/badge/Deploying%20with-Inertia-blue.svg)](https://github.com/ubclaunchpad/inertia)

## Getting Started
Play the latest version [here](http://bumper.ubclaunchpad.com)!  

## Docker Quickstart
Create a .env file in the root directory containing the following variables:
```
NODE_ENV=$YOUR_ENVIRONMENT
SERVER_URL=$YOUR_SERVER_URL
SERVER_PORT=$YOUR_PORT
```

Install Docker and the Docker Compose toolset. Then run:
```bash
docker-compose -f docker-compose.dev.yml up
```

## Manual Quickstart
Go and Node are required. Ensure that environment variables are set or passed into the run commands for the client and server. 

### Server
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

### Client
```bash
$ cd ./client
$ npm install
$ npm start
```
Play at localhost:8080!