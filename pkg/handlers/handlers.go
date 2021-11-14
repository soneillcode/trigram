package handlers

import (
	"log"
	"net/http"
)

func Respond(res http.ResponseWriter, req *http.Request, status int, data *string) {
	res.WriteHeader(status)
	if data != nil {
		byteData := []byte(*data)
		_, err := res.Write(byteData)
		if err != nil {
			log.Println(err)
		}
	}
}

func Error(res http.ResponseWriter, req *http.Request, status int) {
	switch status {
	case http.StatusNotFound:
		http.Error(res, "Not Found", http.StatusNotFound)
	case http.StatusBadRequest:
		http.Error(res, "Bad Request", http.StatusBadRequest)
	default:
		http.Error(res, "Internal Server Error", http.StatusInternalServerError)
	}
}
