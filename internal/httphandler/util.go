package httphandler

import (
	"encoding/json"
	"net/http"
)

type ErrorResponseSt struct {
	ErrorMessage string `json:"error_message"`
}

func uCheckErr(w http.ResponseWriter, err error) bool {
	if err != nil {
		uSendErr(w, http.StatusBadRequest, err)
		return true
	}
	return false
}

func uSendErr(w http.ResponseWriter, code int, err error) {
	uRespondJson(w, code, ErrorResponseSt{
		ErrorMessage: err.Error(),
	})
}

func uRespondJson(w http.ResponseWriter, code int, obj any) {
	raw, err := json.Marshal(obj)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, _ = w.Write(raw)
}
