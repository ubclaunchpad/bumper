.PHONY: deploy clean

all: bumper

bumper:
	env GOOS=linux GOARCH=amd64 go build -v

deploy:
	echo "deploy..."

clean:
	rm server