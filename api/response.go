package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

func SendResponse(w http.ResponseWriter, resp interface{}) {
	resStr, err := json.Marshal(&resp)
	if err != nil {
		log.Println("JSON操作出错")
		return
	}

	_, _ = io.WriteString(w, string(resStr))
}
//
//func SendErrorResponse(w http.ResponseWriter, errResp defs.StrMessage) {
//	resStr, _ := json.Marshal(&errResp.Msg)
//	_, _ = io.WriteString(w, string(resStr))
//}
//
//func SendNormalResponse(w http.ResponseWriter, resp string, sc int) {
//	w.WriteHeader(sc)
//	_, _ = io.WriteString(w, resp)
//}
