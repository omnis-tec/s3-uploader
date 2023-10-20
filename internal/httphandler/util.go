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
		w.WriteHeader(http.StatusBadRequest)
		uRespondJson(w, http.StatusBadRequest, ErrorResponseSt{
			ErrorMessage: err.Error(),
		})
		return true
	}
	return false
}

func uRespondJson(w http.ResponseWriter, code int, obj any) {
	raw, err := json.Marshal(obj)
	if err != nil {
		w.WriteHeader(code)
		_, _ = w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(raw)
}
