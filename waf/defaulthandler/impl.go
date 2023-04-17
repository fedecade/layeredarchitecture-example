package defaulthandler

import (
	"encoding/json"
	"net/http"

	"example.layeredarch/logger"
)

type DefaultHandler struct{}

func New() *DefaultHandler {
	return &DefaultHandler{}
}

func (d *DefaultHandler) ResponseJson(
	obj any,
	status int,
	w http.ResponseWriter,
) {
	body, err := json.Marshal(obj)
	if err != nil {
		logger.Error(err)
		d.ResponseError(http.StatusInternalServerError, w, err)
	}

	w.WriteHeader(status)
	w.Header().Add("Content-Type", "application/json")
	w.Header().Add("Content-Length", string(rune(len(body))))
	w.Write(body)
}

func (d *DefaultHandler) ResponseError(
	code int,
	w http.ResponseWriter,
	err error,
) {
	http.Error(w, http.StatusText(code), code)
	body := []byte(err.Error())
	w.Header().Add("Content-Type", "text/plain")
	w.Header().Add("Content-Length", string(rune(len(body))))
	w.Write(body)
}

func (d *DefaultHandler) ResponseEmpty(
	code int,
	w http.ResponseWriter,
) {
	w.WriteHeader(code)
}
