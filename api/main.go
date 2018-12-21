package main

import (
	"Study/video_server/api/auth"
	h "Study/video_server/api/handlers"
	"Study/video_server/api/session"
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
	auth.ValidateUserSession(r)

	m.r.ServeHTTP(w, r)
}

func RegisterHandlers() *httprouter.Router {
	router := httprouter.New()

	router.POST("/user", h.CreateUser)

	router.POST("/user/:username", h.Login)

	router.GET("/user/:username", h.GetUserInfo)

	router.POST("/user/:username/videos", h.AddNewVideo)

	router.GET("/user/:username/videos", h.ListAllVideos)

	router.DELETE("/user/:username/videos/:vid-id", h.DeleteVideo)

	router.POST("/videos/:vid-id/comments", h.PostComment)

	router.GET("/videos/:vid-id/comments", h.ShowComments)

	return router
}

func Prepare() {
	session.LoadSessionsFromDB()
}

func main() {
	Prepare()
	r := RegisterHandlers()
	mh := NewMiddleWareHandler(r)
	http.ListenAndServe(":8000", mh)
}

