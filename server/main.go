package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/nats-io/nats-server/v2/server"
	"github.com/nats-io/nats.go"
)

const HTTPPort = 8080
const NATSPort = 4223

func webhookHandler(nc *nats.Conn) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {

		// just make sure this is an http post
		if req.Method != http.MethodPost {
			http.Error(w, "Invalid Invocation", http.StatusNotFound)
			return
		}

		body, err := ioutil.ReadAll(req.Body)

		if err != nil {
			http.Error(w, "internal error: "+err.Error(), http.StatusInternalServerError)
			log.Println("i/o error: " + err.Error())
			return
		}

		err = nc.Publish("webHook", body)
		log.Print(string(body))

		if err != nil {
			http.Error(w, "internal error: "+err.Error(), http.StatusInternalServerError)
			log.Println("i/o error: " + err.Error())
			return
		}

	}
}

// create and run a NATS server
// returning a reference to a NATS connection to the embedded NATS server
func setupNATs() *nats.Conn {

	// enable websockets
	// and tell the NATS server to allow connections w/o needing TLS
	// TODO: setup authentication potentially??
	opts := &server.Options{}

	opts.Websocket.Host = "localhost"
	opts.Websocket.Port = 4223
	opts.Websocket.NoTLS = true

	natsAddress := "ws://localhost:4223"

	s, err := server.NewServer(opts)
	if err != nil {
		server.PrintAndDie(fmt.Sprintf("%s", err))
	}

	// start the NATS server
	go func() {

		log.Println("NATS server listening via websocket on", natsAddress)

		if err := server.Run(s); err != nil {
			server.PrintAndDie(err.Error())
		}

		s.WaitForShutdown()

	}()

	fmt.Println("Waiting for NATS server to be ready to accept connections...")
	if !s.ReadyForConnections(10 * time.Second) {
		log.Fatal("NATS server never became ready for connections")
	}

	// connect to our embedded NATS server via nats protocol
	nc, err := nats.Connect("nats://localhost:4222")
	if err != nil {
		log.Fatal("could not connect to embedded NATS server", err)
	}

	return nc

}

// setup HTTP server
func setupHTTP(nc *nats.Conn) {

	port := os.Getenv("PORT")
	if port == "" {
		port = strconv.FormatInt(HTTPPort, 10)
	}

	// setup our endpoint for incoming webhook POSTs
	mux := http.NewServeMux()
	mux.HandleFunc("/webhook/", webhookHandler(nc))

	log.Printf("HTTP server listening via http on http://localhost:%s/webhook\n", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))

}

func main() {
	setupHTTP(setupNATs())
}
