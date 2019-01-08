package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type middleWareHandler struct {
	r *httprouter.Router
}

func NewMiddleWareHandler(r *httprouter.Router) http.Handler {
	m := middleWareHandler{}
	m.r = r
	return m
}

func (m middleWareHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//check session
	ValidateUserSession(r)

	m.r.ServeHTTP(w, r)
}

func RegisterHandlers() *httprouter.Router {
	router := httprouter.New()

	router.POST("/login", BackgroundLogin)
	router.POST("/changepassword", ChangePassword)
	router.GET("/bookmanagerlist", BookManagerList)
	router.POST("/addbook", AddBook)
	router.POST("/updatebook", UpdateBook)
	router.POST("/findbarcodelist", FindBarcodeList)
	router.POST("/changebarcodestatus", ChangeBarcodeStatus)
	router.POST("/addbarcode", AddBarcode)
	router.POST("/deletebarcode", DeleteBarcode)
	router.POST("/exportBarcode.xlsx", ExportBarcode)
	router.POST("/barcodecount", BarcodeCount)


	router.POST("/app_user_login", AppUserLogin)
	router.POST("/app_change_user_name", AppChangeUserName)
	router.POST("/app_unlocked", AppUserUnlocked)
	router.POST("/app_upload_video/:accesstoken&bookid", AppUploadVideo)
	router.POST("/app_get_video_url", AppGetVideoUrl)

	return router
}

func Prepare() {
	LoadSessionsFromDB()
}

func main() {
	Prepare()

	r := RegisterHandlers()
	mh := NewMiddleWareHandler(r)
	_ = http.ListenAndServe(":8000", mh)
}

