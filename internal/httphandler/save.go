package httphandler

import (
	"net/http"
)

type saveRepSt struct {
	Id  string `json:"id"`
	Url string `json:"url"`
}

func (h *handlerSt) Save(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")

	result := make([]saveRepSt, 0, 10)

	if contentType == "multipart/form-data" {
		err := r.ParseMultipartForm(32 << 20)
		if uCheckErr(w, err) {
			return
		}

		for _, fHeaders := range r.MultipartForm.File {
			for _, fHeader := range fHeaders {
				f, err := fHeader.Open()
				if uCheckErr(w, err) {
					return
				}
				defer f.Close()

				repObj, err := h.cr.Save(f, fHeader.Size, fHeader.Header.Get("Content-Type"))
				if uCheckErr(w, err) {
					return
				}

				result = append(result, saveRepSt{
					Id:  repObj.Id,
					Url: repObj.Url,
				})
			}
		}
	} else {
		repObj, err := h.cr.Save(r.Body, r.ContentLength, contentType)
		if uCheckErr(w, err) {
			return
		}

		result = append(result, saveRepSt{
			Id:  repObj.Id,
			Url: repObj.Url,
		})
	}

	uRespondJson(w, http.StatusOK, result)
}
