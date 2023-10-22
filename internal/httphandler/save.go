package httphandler

import (
	"net/http"
	"strconv"
)

type saveRepSt struct {
	Id  string `json:"id"`
	Url string `json:"url"`
}

func (h *handlerSt) Save(w http.ResponseWriter, r *http.Request) {
	var bodySize int64

	if cl := r.Header.Get("Content-Length"); cl != "" {
		bodySize, _ = strconv.ParseInt(cl, 10, 64)
	}

	repObj, err := h.cr.Save(r.Body, bodySize, r.Header.Get("Content-Type"))
	if uCheckErr(w, err) {
		return
	}

	uRespondJson(w, http.StatusOK, saveRepSt{
		Id:  repObj.Id,
		Url: repObj.Url,
	})
}
