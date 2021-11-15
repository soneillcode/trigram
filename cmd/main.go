package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/soneillcode/trigram/pkg/nlp"
)

func main() {
	var port string
	flag.StringVar(&port, "port", "8080", "the port the service is listening on")
	flag.Parse()

	logger := log.New(os.Stdout, "nlp: ", log.Ldate|log.Ltime|log.Lshortfile)
	random := rand.New(rand.NewSource(time.Now().UnixNano()))

	const defaultNumberOfWords = 200
	nlpService := nlp.NewService(random, defaultNumberOfWords)
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
