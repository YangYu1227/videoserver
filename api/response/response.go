package response

import (
	"Study/video_server/api/defs"
	"encoding/json"
	"io"
	"net/http"
)

func SendErrorResponse(w http.ResponseWriter, errResp defs.ErrorResponse) {
	w.WriteHeader(errResp.HttpStateCode)
	resStr, _ := json.Marshal(&errResp.Error)
	_, _ = io.WriteString(w, string(resStr))
}

func SendNormalResponse(w http.ResponseWriter, resp string, sc int) {
	w.WriteHeader(sc)
	_, _ = io.WriteString(w, resp)
}
