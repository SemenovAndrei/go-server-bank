package main

import (
	"github.com/i-hit/go-server-bank.git/cmd/bank/app"
	"github.com/i-hit/go-server-bank.git/pkg/card"
	"log"
	"net"
	"net/http"
	"os"
)

const defaultPort = "9999"
const defaultHost = "0.0.0.0"

func main() {
	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = defaultPort
	}

	host, ok := os.LookupEnv("HOST")
	if !ok {
		host = defaultHost
	}

	log.Println(host)
	log.Println(port)

	if err := execute(net.JoinHostPort(host, port)); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

func execute(adrdr string) (err error) {
	cardSvc := card.NewService()
	mux := http.NewServeMux()
	application := app.NewServer(cardSvc, mux)

	_, _ = cardSvc.Add("12", "basic", "Visa")
	_, _ = cardSvc.Add("12", "virtual", "Visa")

	application.Init()

	server := &http.Server{
		Addr: adrdr,
		Handler: application,
	}
	return server.ListenAndServe()
}
