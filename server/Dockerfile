FROM golang:alpine

WORKDIR /app
ENV SRC_DIR=/go/src/github.com/ubclaunchpad/bumper/server

# Copy code over from source directory and build server
ADD . $SRC_DIR
RUN cd $SRC_DIR; go build -o server; cp server /app

# Expose ports and run application
EXPOSE 9090
ENTRYPOINT ./server
