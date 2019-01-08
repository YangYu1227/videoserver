package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func RegisterHandler() *httprouter.Router {
	router := httprouter.New()

	router.GET("/", LoginHandler)
	router.GET("/bookmanager", BookManagerHandler)
	router.GET("/bookbarcode/:book_id", BookBarcodeHandler)
	router.GET("/bookbarcodecount/:book_id", BookBarcodeCount)

	router.POST("/api", APIHandler)


	router.ServeFiles("/statics/*filepath", http.Dir("./templates"))

	return router
}

func main() {
	r := RegisterHandler()
	_ = http.ListenAndServe(":8080", r)
}
