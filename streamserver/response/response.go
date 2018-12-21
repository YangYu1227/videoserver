package response

import (
	"io"
	"net/http"
)

func SendErrorResponse(w http.ResponseWriter, sc int, errMsg string) {
	w.WriteHeader(sc)
	_, _ = io.WriteString(w, errMsg)
}