package nlp

import (
	"io/ioutil"
	"log"
	"net/http"

	"example.com/todo/pkg/handlers"
)

// Handler encapsulates http functionality and concerns. It stores a service which handles the actual functionality.
type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

// Learn handles http POST requests of text data and adds it to the existing body of data.
func (h *Handler) Learn(res http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		handlers.Error(res, req, http.StatusBadRequest)
		return
	}

	buf, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Printf("failed to read request body: %v", err)
		handlers.Error(res, req, http.StatusBadRequest)
		return
	}
	if len(buf) == 0 {
		log.Printf("request body is empty")
		handlers.Error(res, req, http.StatusBadRequest)
		return
	}

	h.service.Learn(string(buf))

	handlers.Respond(res, req, http.StatusOK, nil)
}

// Generate handles http GET requests, randomly generating a new sample of text using the body of data and returning it.
func (h *Handler) Generate(res http.ResponseWriter, req *http.Request) {
	if req.Method != "GET" {
		handlers.Error(res, req, http.StatusBadRequest)
		return
	}

	data, err := h.service.Generate()
	if err != nil {
		log.Printf("service failed to handle 'generate': %v", err)
		handlers.Error(res, req, http.StatusInternalServerError)
		return
	}
	handlers.Respond(res, req, http.StatusOK, data)
}
