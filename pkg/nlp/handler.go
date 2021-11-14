package nlp

import (
	"io/ioutil"
	"log"
	"net/http"
)

// Handler encapsulates http functionality and concerns. It stores a service which handles the actual functionality.
type Handler struct {
	service *Service
	logger  *log.Logger
}

// NewHandler returns a new instance of a Handler
func NewHandler(service *Service, logger *log.Logger) *Handler {
	return &Handler{
		service: service,
		logger:  logger,
	}
}

// Learn handles http POST requests. It takes text data from the request and adds it to the existing data.
func (h *Handler) Learn(res http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		h.respondWithError(res, http.StatusBadRequest)
		return
	}

	buf, err := ioutil.ReadAll(req.Body)
	if err != nil {
		h.logger.Printf("failed to read request body: %v", err)
		h.respondWithError(res, http.StatusBadRequest)
		return
	}
	if len(buf) == 0 {
		h.logger.Printf("request body is empty")
		h.respondWithError(res, http.StatusBadRequest)
		return
	}

	h.service.Learn(string(buf))
	h.respond(res, http.StatusOK, nil)
}

// Generate handles http GET requests. It randomly generates a new sample of text using the stored data and returns it.
func (h *Handler) Generate(res http.ResponseWriter, req *http.Request) {
	if req.Method != "GET" {
		h.respondWithError(res, http.StatusBadRequest)
		return
	}

	data, err := h.service.Generate()
	if err != nil {
		h.logger.Printf("service failed to handle 'generate': %v", err)
		h.respondWithError(res, http.StatusInternalServerError)
		return
	}
	h.respond(res, http.StatusOK, data)
}

func (h *Handler) respond(res http.ResponseWriter, status int, data *string) {
	res.WriteHeader(status)
	if data != nil {
		byteData := []byte(*data)
		_, err := res.Write(byteData)
		if err != nil {
			h.logger.Println(err)
		}
	}
}

func (h *Handler) respondWithError(res http.ResponseWriter, status int) {
	switch status {
	case http.StatusNotFound:
		http.Error(res, "Not Found", http.StatusNotFound)
	case http.StatusBadRequest:
		http.Error(res, "Bad Request", http.StatusBadRequest)
	default:
		http.Error(res, "Internal Server Error", http.StatusInternalServerError)
	}
}
