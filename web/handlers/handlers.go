package handlers

import (
	"Study/video_server/web/client"
	"Study/video_server/web/defs"
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

func HomeHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	cname, err1 := r.Cookie("username")
	sid, err2 := r.Cookie("session")

	if err1 != nil || err2 != nil {
		p := &HomePage{Name: "yy"}
		t, e := template.ParseFiles("./templates/home.html")
		if e != nil {
			log.Printf("Parsing template home.html error: %s", e)
			return
		}

		_ = t.Execute(w, p)
		return
	}

	if len(cname.Value) != 0 && len(sid.Value) != 0 {
		http.Redirect(w, r, "/userhome", http.StatusFound)
		return
	}
}

func UserHomeHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	cname, err1 := r.Cookie("username")
	_, err2 := r.Cookie("session")

	if err1 != nil || err2 != nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	fname := r.FormValue("username")

	var p *UserPage
	if len(cname.Value) != 0 {
		p = &UserPage{Name: cname.Value}
	} else if len(fname) != 0 {
		p = &UserPage{Name: fname}
	}

	t, e := template.ParseFiles("./templates/userhome.html")
	if e != nil {
		log.Printf("Parsing userhome.html error: %s", e)
		return
	}

	_ = t.Execute(w, p)
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
	client.Request(apibody, w, r)
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
