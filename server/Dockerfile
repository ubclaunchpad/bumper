FROM golang:alpine

WORKDIR /app
ENV SRC_DIR=/go/src/github.com/ubclaunchpad/bumper/server

# Copy code over from source directory and build server
RUN apk add --update --no-cache git
ADD . $SRC_DIR
WORKDIR $SRC_DIR
RUN if [ ! -d "vendor" ]; then \
    go get -u github.com/golang/dep/cmd/dep; \
    dep ensure; \
    fi
RUN go build -o server; cp server /app

# Expose ports and run application
EXPOSE 9090
ENTRYPOINT ./server
