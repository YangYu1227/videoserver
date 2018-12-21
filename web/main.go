package main

import (
	"Study/video_server/web/handlers"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func RegisterHandler() *httprouter.Router {
	router := httprouter.New()

	router.GET("/", handlers.HomeHandler)
	router.POST("/", handlers.HomeHandler)
	router.GET("/userhome", handlers.UserHomeHandler)
	router.POST("/userhome", handlers.UserHomeHandler)

	router.POST("/api", handlers.APIHandler)

	router.GET("/videos/:vid-id", handlers.ProxyVideoHandler)
	router.POST("/upload/:vid-id", handlers.ProxyUploadHandler)

	router.ServeFiles("/statics/*filepath", http.Dir("./templates"))

	return router
}

func main() {
	r := RegisterHandler()
	http.ListenAndServe(":8080", r)
}
