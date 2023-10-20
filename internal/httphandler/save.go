package httphandler

import (
	"net/http"
)

type saveRepSt struct {
	Id  string `json:"id"`
	Url string `json:"url"`
}

func (h *handlerSt) Save(w http.ResponseWriter, r *http.Request) {
	repObj, err := h.cr.Save(r.Body, r.Header.Get("Content-Type"))
	if uCheckErr(w, err) {
		return
	}

	uRespondJson(w, http.StatusOK, saveRepSt{
		Id:  repObj.Id,
		Url: repObj.Url,
	})
}
