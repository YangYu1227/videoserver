package handlers

import (
	"Study/video_server/scheduler/dbops"
	"Study/video_server/scheduler/response"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func VidDelRecHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	vid := p.ByName("vid-id")

	if len(vid) == 0 {
		response.SendResponse(w, 400, "video id should not be empty")
		return
	}

	err := dbops.AddVideoDeletionRecord(vid)
	if err != nil {
		response.SendResponse(w, 500, "Internal server error")
		return
	}

	response.SendResponse(w, 200, "")
	return
}
