package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/soneillcode/trigram/pkg/nlp"
)

func main() {
	var port string
	flag.StringVar(&port, "port", "8080", "the port the service is listening on")
	flag.Parse()

	logger := log.New(os.Stdout, "nlp: ", log.Ldate|log.Ltime|log.Lshortfile)

	nlpService := nlp.NewService()
	nlpHandler := nlp.NewHandler(nlpService, logger)

	http.HandleFunc("/learn", nlpHandler.Learn)
	http.HandleFunc("/generate", nlpHandler.Generate)

	logger.Printf("starting server on port: %s\n", port)
	server := http.Server{
		Addr:     fmt.Sprintf(":%s", port),
		ErrorLog: logger,
	}
	err := server.ListenAndServe()
	if err != nil {
		logger.Fatalf("failed to listen and serve: %s", err)
	}
}
