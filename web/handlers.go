package main

import (
	"Study/book_server/web/defs"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"yy/config"
)

type HomePage struct {
	Name string
}

type UserPage struct {
	Name string
}

func LoginHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	t, err := template.ParseFiles("./templates/login/login.html")
	if err != nil {
		log.Printf("Parsing template home.html error: %s", err)
	}
	_ = t.Execute(w, nil)
}

func BookManagerHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	t, err := template.ParseFiles("./templates/bookmanger/bookmanger_list.html")
	if err != nil {
		log.Printf("Parsing template home.html error: %s", err)
	}
	_ = t.Execute(w, nil)
}

func BookBarcodeHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	bookId := ps.ByName("book_id")
	log.Printf("获取到的bookid: %s", bookId)

	t, err := template.ParseFiles("./templates/decode/decode_list.html")
	if err != nil {
		log.Printf("Parsing template home.html error: %s", err)
	}
	_ = t.Execute(w, nil)
}

func BookBarcodeCount(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	bookId := ps.ByName("book_id")
	log.Printf("获取到的bookid: %s", bookId)

	t, err := template.ParseFiles("./templates/decode/decode_count.html")
	if err != nil {
		log.Printf("Parsing template home.html error: %s", err)
	}
	_ = t.Execute(w, nil)
}

func APIHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if r.Method != http.MethodPost {
		re, _ := json.Marshal(defs.ErrorRequestNotRecognized)
		_, _ = io.WriteString(w, string(re))
		return
	}

	res, _ := ioutil.ReadAll(r.Body)
	apibody := &defs.ApiBody{}
	if err := json.Unmarshal(res, apibody); err != nil {

		re, _ := json.Marshal(defs.ErrorRequestBodyParseFailed)
		_, _ = io.WriteString(w, string(re))
		return
	}
	log.Println(apibody.Url)
	Request(apibody, w, r)
	defer r.Body.Close()
}

func ProxyVideoHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	u, _ := url.Parse("http://" + config.GetLbAddr() + ":9000/")
	proxy := httputil.NewSingleHostReverseProxy(u)
	proxy.ServeHTTP(w, r)
}

func ProxyUploadHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	log.Println("Web的ProxyUploadHandler方法被调用")
	u, _ := url.Parse("http://" + config.GetLbAddr() + ":9000/")
	log.Println(u)
	proxy := httputil.NewSingleHostReverseProxy(u)
	proxy.ServeHTTP(w, r)
}
