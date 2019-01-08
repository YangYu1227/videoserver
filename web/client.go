package main

import (
	"Study/book_server/web/defs"
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"yy/config"
)

var httpClient *http.Client

func init() {
	httpClient = &http.Client{}
}

func Request(b *defs.ApiBody, w http.ResponseWriter, r *http.Request) {
	var resp *http.Response
	var err error

	u, _ := url.Parse(b.Url)
	u.Host = config.GetLbAddr() + ":" + u.Port()
	newUrl := u.String()

	log.Printf("拼接后的地址：%s", newUrl)

	switch b.Method {
	case http.MethodGet:
		log.Printf("调用API的GET方法，地址：%s", newUrl)
		req, _ := http.NewRequest("GET", newUrl, nil)
		req.Header = r.Header
		resp, err = httpClient.Do(req)
		if err != nil {
			log.Println(err)
			return
		}
		normalResponse(w, resp)
	case http.MethodPost:
		req, _ := http.NewRequest("POST", newUrl, bytes.NewBuffer([]byte(b.ReqBody)))
		req.Header = r.Header
		resp, err = httpClient.Do(req)
		if err != nil {
			log.Println(err)
			return
		}
		normalResponse(w, resp)
	case http.MethodDelete:
		req, _ := http.NewRequest("DELETE", newUrl, nil)
		req.Header = r.Header
		resp, err = httpClient.Do(req)
		if err != nil {
			log.Println(err)
			return
		}
		normalResponse(w, resp)
	default:
		w.WriteHeader(http.StatusBadRequest)
		_, _ = io.WriteString(w, "Bad api request")
		return
	}
}

func normalResponse(w http.ResponseWriter, r *http.Response) {
	res, err := ioutil.ReadAll(r.Body)
	if err != nil {
		re, _ := json.Marshal(defs.ErrorInternalFaults)
		w.WriteHeader(500)
		_, _ = io.WriteString(w, string(re))
		return
	}

	w.WriteHeader(r.StatusCode)
	_, _ = io.WriteString(w, string(res))
}
