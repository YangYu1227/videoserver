package main

import (
	"Study/video_server/scheduler/handlers"
	"Study/video_server/scheduler/taskrunner"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func RegisterHandlers() *httprouter.Router {
	router := httprouter.New()
	router.GET("/video-delete-record/:vid-id", handlers.VidDelRecHandler)

	return router
}

func main() {
	go taskrunner.Start()
	r := RegisterHandlers()
	http.ListenAndServe(":9001", r)
}
