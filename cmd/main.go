package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"example.com/todo/pkg/nlp"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("%s\n", err)
	}
}

func run() error {
	var port string
	flag.StringVar(&port, "port", "8080", "the port the service is listening on")
	flag.Parse()

	nlpService := nlp.NewService()
	nlpHandler := nlp.NewHandler(nlpService)

	http.HandleFunc("/learn", nlpHandler.Learn)
	http.HandleFunc("/generate", nlpHandler.Generate)

	log.Printf("starting server on port: %s\n", port)
	err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
	if err != nil {
		return fmt.Errorf("failed to listen and serve: %w", err)
	}

	return nil
}
