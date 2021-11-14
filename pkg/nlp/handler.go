package nlp

import (
	"io/ioutil"
	"log"
	"net/http"

	"example.com/todo/pkg/routing"
)

// Handler encapsulates http functionality and concerns. It stores a service which handles the actual functionality.
type Handler struct {
	service              *Service
	defaultNumberOfWords int
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service:              service,
		defaultNumberOfWords: 400, // consider overriding this from a request param
	}
}

// Learn handles http POST requests of text data and adds it to the existing body of data.
func (h *Handler) Learn(res http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		routing.Error(res, req, http.StatusBadRequest)
		return
	}

	buf, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Printf("failed to read request body: %v", err)
		routing.Error(res, req, http.StatusBadRequest)
		return
	}
	if len(buf) == 0 {
		log.Printf("request body is empty")
		routing.Error(res, req, http.StatusBadRequest)
		return
	}

	err = h.service.Learn(string(buf))
	if err != nil {
		log.Printf("service failed to handle 'learn': %v", err)
		routing.Error(res, req, http.StatusInternalServerError)
		return
	}

	routing.Respond(res, req, http.StatusOK, nil)
}

// Generate handles http GET requests, randomly generating a new sample of text using the body of data and returning it.
func (h *Handler) Generate(res http.ResponseWriter, req *http.Request) {
	if req.Method != "GET" {
		routing.Error(res, req, http.StatusBadRequest)
		return
	}

	data, err := h.service.Generate(h.defaultNumberOfWords)
	if err != nil {
		log.Printf("service failed to handle 'generate': %v", err)
		routing.Error(res, req, http.StatusInternalServerError)
		return
	}
	routing.Respond(res, req, http.StatusOK, data)
}
