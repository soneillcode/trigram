package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"example.com/todo/pkg/nlp"
	"example.com/todo/pkg/routing"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run() error {

	var port string
	flag.StringVar(&port, "port", "8080", "the port the service is listening on")
	flag.Parse()

	nlpService := nlp.NewService()
	nlpHandler := nlp.NewHandler(nlpService)

	servicesRouter := routing.NewRouter()

	servicesRouter.AddRoute("/learn", "POST", nlpHandler.Learn)
	servicesRouter.AddRoute("/generate", "GET", nlpHandler.Generate)

	// start server
	log.Printf("starting server on port: %s\n", port)
	err := http.ListenAndServe(fmt.Sprintf(":%s", port), servicesRouter)
	if err != nil {
		return fmt.Errorf("failed to listen and serve: %w", err)
	}

	return nil
}
