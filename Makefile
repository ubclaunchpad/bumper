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

# Generate self-signed key and certificate
.PHONY: https
https:
	openssl req -x509 -newkey rsa:4096 -keyout config/key.pem -out config/cert.pem -days 365 -nodes

# Run NGINX container in daemon mode
# KNOWN ISSUE: Docker for Mac/Windows runs in VM --network host will not work as expected
# Replace --network host with -p 80:80 and -p 443:443 when in development
.PHONY: nginx
nginx:
	docker run -d --rm \
		--name bumper_nginx \
		-v $(PWD)/config/nginx.conf:/etc/nginx/nginx.conf \
		-v $(PWD)/config/cert.pem:/etc/nginx/ssl/nginx.crt \
		-v $(PWD)/config/key.pem:/etc/nginx/ssl/nginx.key \
		--network host \
		nginx:alpine

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

# Starts a file server in the web/ directory
.PHONY: web
web:
	(cd ./web ; python -m SimpleHTTPServer 8000)

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