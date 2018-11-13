# Installs project dependencies
.PHONY: deps
deps:
	go get github.com/codegangsta/gin
	(cd server ; go get -u github.com/golang/dep/cmd/dep ; dep ensure)
	(cd client ; npm install )

# Build and minify client
.PHONY: bundle
bundle:
	(cd ./client ; npm run build)

# Deploy to GH Pages
.PHONY: deploy
deploy:
	make bundle
	git subtree push --prefix client/public origin gh-pages

# Build and run Bumper in daemon mode
.PHONY: bumper
DATABASE_URL=https://bumperdevdb.firebaseio.com
SERVER_PORT=9090
bumper:
	docker stop bumper 2>&1 || true
	docker build -t bumper .
	docker run -d --rm \
		--name bumper \
		-e DATABASE_URL=$(DATABASE_URL) \
		-e PORT=$(SERVER_PORT) \
		-v $(PWD)/client/public:/app/build \
		-p 80:$(SERVER_PORT) \
		bumper

# Starts the client (dev server on port 8080)
.PHONY: client
client:
	(cd ./client ; npm start)

# Starts the server (exposed on port 9090)
.PHONY: server
server:
	(cd ./server ; DATABASE_URL=$(DATABASE_URL) PORT=8081 gin -p $(SERVER_PORT) -a 8081 -i run main.go)

# Runs unit tests
.PHONY: test
test:
	(cd ./server ; go test -race ./...)