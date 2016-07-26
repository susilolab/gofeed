package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/goincremental/negroni-sessions"
	"github.com/unrolled/render"
)

func userIndex(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	session := sessions.GetSession(req)
	user := session.Get("user").([]byte)

	renderOpts.Layout = "admin"
	rnd = render.New(renderOpts)
	response := NewResponse()
	response.Title = "Index"
	var u User
	if err := json.Unmarshal(user, &u); err != nil {
		log.Println(err)
	}
	response.Data = u

	rnd.HTML(w, http.StatusOK, "user/index", response)
}

func userUpdate(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
}

func userDelete(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
}

func userView(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
}

func userMiddleware(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	session := sessions.GetSession(req)
	id := session.Get("user_id")

	if id == nil {
		http.Redirect(w, req, "/site/unauthorized", 301)
	} else {
		next(w, req)
	}
}
