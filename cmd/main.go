package main

import (
	"flag"
	"log"
	"net"
	"net/http"
	"strconv"

	"github.com/google/go-safeweb/safehttp"

	"github.com/aiit2022-pbl-okuhara/incident-training/secure"
	"github.com/aiit2022-pbl-okuhara/incident-training/server"
	"github.com/aiit2022-pbl-okuhara/incident-training/storage"
)

var (
	port = flag.Int("port", 8080, "Port for the HTTP server")
	dev  = flag.Bool("dev", false, "Run in development mode")
)

func main() {
	log.SetFlags(log.Flags() | log.Lshortfile)
	flag.Parse()
	safehttp.UseLocalDev()
	if *dev {
		safehttp.UseLocalDev()
	}
	db, err := storage.NewDB()
	if err != nil {
		// TODO: 適切に error 処理を行う
		log.Fatal(err)
	}

	addr := net.JoinHostPort("0.0.0.0", strconv.Itoa(*port))
	mux := secure.NewMuxConfig(db, addr).Mux()
	server.Load(db, mux)

	log.Printf("Listening on %q", addr)

	// TODO: 適切に error 処理を行う
	log.Fatal(http.ListenAndServe(addr, mux))
}
