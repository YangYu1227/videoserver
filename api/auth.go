package main

import (
	"net/http"
)

var HEADER_FIELD_SESSION = "X-Session-Id"
var HEADER_FIELD_UNAME = "X-User-Name"

func ValidateUserSession(r *http.Request) bool {
	sid := r.Header.Get(HEADER_FIELD_SESSION)
	if len(sid) == 0 {
		return false
	}

	uname, ok := IsSessionExpired(sid)
	if ok {
		return false
	}

	r.Header.Add(HEADER_FIELD_UNAME, uname)
	return true
}


func ValidateUser(w http.ResponseWriter, r *http.Request) bool {
	uname := r.Header.Get(HEADER_FIELD_UNAME)
	//log.Printf("uname:%s\n", uname)
	if len(uname) == 0 {
		return false
	}
	return true
}