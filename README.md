# bumper
bump other players and objects off to stay alive

[![Go Report Card](https://goreportcard.com/badge/github.com/ubclaunchpad/bumper)](https://goreportcard.com/report/github.com/ubclaunchpad/bumper) [![Deployed with Inertia](https://img.shields.io/badge/Deploying%20with-Inertia-blue.svg)](https://github.com/ubclaunchpad/inertia)

## Getting Started
Play the latest version [here](https://bumper.ubclaunchpad.com)!  

## Development
### Running the server
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
```bash
$ cd ./client
$ npm install
$ npm start
```
Play at localhost:8080!
