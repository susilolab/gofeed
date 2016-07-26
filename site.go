package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/goincremental/negroni-sessions"
	"github.com/gorilla/schema"
	"github.com/unrolled/render"
)

func siteView(w http.ResponseWriter, req *http.Request) {
}

func siteLogin(w http.ResponseWriter, req *http.Request) {
	session := sessions.GetSession(req)
	err := session.Get("error")

	rnd = render.New(renderOpts)
	response := NewResponse()
	response.Title = "Login"
	response.Error = 0
	response.Message = map[string]interface{}{
		"username":   "",
		"password":   "",
		"user_error": 0,
		"pass_error": 0,
	}
	response.Data = UserLogin{}

	if err != nil {
		var u UserLogin
		sUser := session.Get("user")
		if sUser != nil {
			if errUser := json.Unmarshal(sUser.([]byte), &u); errUser != nil {
				log.Println(errUser)
			}
		}
		response.Error = 1
		response.Message = err
		response.Data = u
		session.Delete("error")
		session.Delete("user")
	}
	rnd.HTML(w, http.StatusOK, "site/login", response)
}

func siteLoginPost(w http.ResponseWriter, req *http.Request) {
	session := sessions.GetSession(req)
	decoder := schema.NewDecoder()
	err := req.ParseForm()
	if err != nil {
		log.Println(err)
	}

	userLogin := &UserLogin{}
	err = decoder.Decode(userLogin, req.PostForm)
	if err != nil {
		log.Println(err)
	}

	var user User
	validate, errMsg := userLogin.validate()
	ul, _ := userLogin.toJson()
	session.Set("user", ul)
	if validate {
		pass := hashPassword(userLogin.Password)
		if db.Table("users").Where(&User{Username: userLogin.Username,
			PasswordHash: pass}).
			First(&user).RecordNotFound() {

			session.Set("error", map[string]interface{}{
				"password":   "Username/Password tidak benar.",
				"user_error": 0,
				"pass_error": 1,
			})
			http.Redirect(w, req, "/site/login", 301)
		} else {
			session.Set("user_id", user.ID)
			u, _ := user.toJson()
			session.Set("user", u)
			http.Redirect(w, req, "/user/index", 301)
		}

	} else {
		session.Set("error", errMsg)
		http.Redirect(w, req, "/site/login", 301)
	}
}

func siteLogout(w http.ResponseWriter, req *http.Request) {
	session := sessions.GetSession(req)
	session.Delete("user_id")
	session.Delete("user")
	http.Redirect(w, req, "/", 301)
}

func siteSignup(w http.ResponseWriter, req *http.Request) {
	session := sessions.GetSession(req)
	err := session.Get("error")

	rnd = render.New(renderOpts)
	response := NewResponse()
	response.Title = "Signup"
	response.Data = User{}
	if err != nil {
		var u User
		sUser := session.Get("user")
		if sUser != nil {
			if errUser := json.Unmarshal(sUser.([]byte), &u); errUser != nil {
				log.Println(errUser)
			}
		}
		response.Data = u
		response.Error = 1
		response.Message = err
		session.Delete("error")
		session.Delete("user")
	}
	rnd.HTML(w, http.StatusOK, "site/signup", response)
}

func siteSignupPost(w http.ResponseWriter, req *http.Request) {
	session := sessions.GetSession(req)
	decoder := schema.NewDecoder()
	err := req.ParseForm()
	if err != nil {
		log.Println(err)
	}

	user := &User{}
	err = decoder.Decode(user, req.PostForm)
	if err != nil {
		log.Println(err)
	}

	validate, errMsg := user.validate()
	u, _ := user.toJson()
	session.Set("user", u)
	if validate {
		rnd = render.New(renderOpts)
		response := NewResponse()
		response.Title = "Signup"
		response.Data = User{}

		user.Status = 1
		user.PasswordHash = hashPassword(user.PasswordHash)
		user.CreatedAt = time.Now()
		user.UpdatedAt = time.Now()
		if errCreate := db.Table("users").Create(&user).Error; errCreate == nil {
			response.Message = map[string]interface{}{
				"error": 0,
				"msg":   "User berhasil ditambah.",
			}
		} else {
			response.Error = 1
			response.Message = map[string]interface{}{
				"error": 1,
				"msg":   "User gagal ditambah.",
			}
		}
		rnd.HTML(w, http.StatusOK, "site/signup", response)
	} else {
		session.Set("error", errMsg)
		http.Redirect(w, req, "/site/signup", 301)
	}
}

func siteResetPassword(w http.ResponseWriter, req *http.Request) {
}

func siteResetPasswordPost(w http.ResponseWriter, req *http.Request) {
}

func siteReqResetPassword(w http.ResponseWriter, req *http.Request) {
}

func siteReqResetPasswordPost(w http.ResponseWriter, req *http.Request) {
}

func siteAbout(w http.ResponseWriter, req *http.Request) {
	rnd = render.New(renderOpts)
	response := NewResponse()
	response.Title = "About"
	rnd.HTML(w, http.StatusOK, "site/about", response)
}

func siteUnauthorized(w http.ResponseWriter, req *http.Request) {
	rnd = render.New(renderOpts)
	response := NewResponse()
	response.Title = "Akses ditolak"
	rnd.HTML(w, http.StatusOK, "site/unauthorized", response)
}
