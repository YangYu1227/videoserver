package handlers

import (
	"Study/video_server/streamserver/defs"
	"Study/video_server/streamserver/response"
	"github.com/julienschmidt/httprouter"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"yy/ossops"
)

func StreamHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	//vid := p.ByName("vid-id")
	//vl := defs.VIDEO_DIR + vid
	//log.Printf("视频地址：%s", vl)
	//video, err := os.Open(vl)
	//if err != nil {
	//	log.Printf("Error when try to open file: %v", err)
	//	response.SendErrorResponse(w, http.StatusInternalServerError, "Internal Error")
	//	return
	//}
	//
	//w.Header().Set("Content-Type", "video/mp4")
	//http.ServeContent(w, r, "", time.Now(), video)
	//
	//defer video.Close()

	log.Println("Entered the streamHandler")
	targetUrl := "http://yy-video-server-oss-xianggang.oss-cn-hongkong.aliyuncs.com/videos/" + p.ByName("vid-id")
	http.Redirect(w, r, targetUrl, 301)
}

func UploadHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	r.Body = http.MaxBytesReader(w, r.Body, defs.MAX_UPLOAD_SIZE)
	if err := r.ParseMultipartForm(defs.MAX_UPLOAD_SIZE); err != nil {
		response.SendErrorResponse(w, http.StatusBadRequest, "File is too big")
		return
	}

	file, _, err := r.FormFile("file") // <from name="file">
	if err != nil {
		response.SendErrorResponse(w, http.StatusInternalServerError, "Internal Error")
		return
	}

	data, err := ioutil.ReadAll(file)
	if err != nil {
		log.Printf("Read file error: %v", err)
		response.SendErrorResponse(w, http.StatusInternalServerError, "Internal Error")
		return
	}

	filename := p.ByName("vid-id")
	err = ioutil.WriteFile(defs.VIDEO_DIR+filename, data, 0666)
	if err != nil {
		log.Printf("Write file error: %v", err)
		response.SendErrorResponse(w, http.StatusInternalServerError, "Internal Error")
		return
	}

	ossfn := "videos/" + filename
	path := "./videos/" + filename
	bn := "yy-video-server-oss-xianggang"
	ret := ossops.UploadToOss(ossfn, path, bn)
	if !ret {
		response.SendErrorResponse(w, http.StatusInternalServerError, "Internal Error")
		return
	}

	_ = os.Remove(path)

	w.WriteHeader(http.StatusCreated)
	_, _ = io.WriteString(w, "Uploaded successfully")
}
