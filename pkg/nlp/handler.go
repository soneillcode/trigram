package nlp

import (
	"io/ioutil"
	"log"
	"net/http"

	"example.com/todo/pkg/routing"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) Learn(res http.ResponseWriter, req *http.Request) {
	buf, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Printf("failed to read request body: %v", err)
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

func (h *Handler) Generate(res http.ResponseWriter, req *http.Request) {
	data, err := h.service.Generate()
	if err != nil {
		log.Printf("service failed to handle 'generate': %v", err)
		routing.Error(res, req, http.StatusInternalServerError)
		return
	}
	routing.Respond(res, req, http.StatusOK, data)
}
