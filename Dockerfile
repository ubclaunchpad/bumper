# Build and minify React client
FROM node:carbon AS client
WORKDIR /client
ADD client .
RUN npm install
RUN npm run build

# Build server
FROM golang:alpine AS server
WORKDIR /app
ENV SRC_DIR=/go/src/github.com/ubclaunchpad/bumper/server
RUN apk add --update --no-cache git
ADD server $SRC_DIR
WORKDIR $SRC_DIR
RUN if [ ! -d "vendor" ]; then \
    go get -u github.com/golang/dep/cmd/dep; \
    dep ensure; \
    fi
RUN go build -o server; cp server /app

# Copy build to final stage
FROM alpine
WORKDIR /app/build
COPY --from=client /client/public/ .
WORKDIR /app
COPY --from=server /app/server .

EXPOSE 80
ENTRYPOINT ./server