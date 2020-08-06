# Build server
FROM golang:alpine AS server
WORKDIR /app
ENV SRC_DIR=/go/src/github.com/ubclaunchpad/bumper/server
RUN apk add --update --no-cache git
ADD server $SRC_DIR
WORKDIR $SRC_DIR
RUN go mod download
RUN go build -o server; cp server /app

# Copy build to final stage
FROM alpine
RUN apk add --update --no-cache ca-certificates
WORKDIR /app
COPY server/service-account.json .
COPY --from=server /app/server .

ENV PORT 80
EXPOSE 80
ENTRYPOINT ./server
