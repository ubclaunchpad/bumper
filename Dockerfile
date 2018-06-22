# Build and minify React client
FROM node:carbon AS client
WORKDIR /client
ADD client/package.json .
RUN npm install
ADD client .
RUN npm run build

# Build server
FROM golang:alpine AS server
WORKDIR /app
ENV SRC_DIR=/go/src/github.com/ubclaunchpad/bumper/server
RUN apk add --update --no-cache git
ADD server $SRC_DIR
WORKDIR $SRC_DIR
RUN go get -u github.com/golang/dep/cmd/dep
RUN dep ensure --vendor-only
RUN go build -o server; cp server /app

# Copy build to final stage
FROM alpine
WORKDIR /app/build
COPY --from=client /client/public/ .
WORKDIR /app
COPY --from=server /app/service-account.json .
COPY --from=server /app/server .

ENV PORT 9090
EXPOSE 9090
ENTRYPOINT ./server
