# webhook2NATs

This is a little server utility that publishes data received via HTTP POSTs as a [NATS](https://nats.io/) message.

This could, for example, be useful if you are dealing with an API that does HTTP POSTs to alert you of something, but you can't get the HTTP Posts on your local development machine.

## Server

This is a Go server that contains an embedded [NATS server](github.com/nats-io/nats-server/v2/server).

Additionally, the [Go NATs client library](github.com/nats-io/nats.go) is used to connect to the embedded NATS server.

Finally the [GO's http package](https://pkg.go.dev/net/http) is used to setup an endpoint that allows HTTP POST requests, which publishes a NATs message with the posted data under the name "webhook".

## Client

This is a node.js server-side client that uses the [NATS client library](github.com/nats-io/nats.go) to connect to the GO server, subscribes to the "webhook" messages, and prints the received data.

(NOTE that obviously the client is just an example for you to use to integrate subscribing to the NATs message into your code.)

## Start Server

```console
cd server
go run main.go
```

## Install Client

```console
cd client
npm install
```

## Run Client

```console
cd client
npm start
```

## Test

```console
curl -v --request POST --url http://localhost:8080/webhook/ --data 'hello'
```

# Example Use Case

- Run this server on a public web server
- Do your local development work on your local laptop, connecting to your server via NATs
