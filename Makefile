# Installs project dependencies
.PHONY: deps
deps:
	go get github.com/codegangsta/gin
	(cd server ; go get -u github.com/golang/dep/cmd/dep ; dep ensure)
	(cd client ; npm install )

# Starts both the server and the client via docker-compose
.PHONY: docker-start
docker-start:
	docker-compose -f docker-compose.dev.yml up

# Starts the client
.PHONY: client
client:
	(cd ./client ; npm start)

# Starts the server
.PHONY: server
server:
	(cd ./server ; PORT=9090 ; go run main.go)
