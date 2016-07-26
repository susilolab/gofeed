/* Gofeed, web application original convert from yiifeed.com. but it is for go news
 *
 * Todo:
 * - Change userMiddleware behavior to prevent assign to router one by one,
 *   just assign to subrouter
 */
package main

import (
	"encoding/gob"
	"encoding/json"
	// "flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"time"

	// "github.com/agus/utils"
	"github.com/codegangsta/negroni"
	_ "github.com/go-sql-driver/mysql"
	"github.com/goincremental/negroni-sessions"
	"github.com/goincremental/negroni-sessions/cookiestore"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/unrolled/render"
	"gopkg.in/tylerb/graceful.v1"
)

var (
	rnd        *render.Render
	renderOpts render.Options
	db         *gorm.DB
	rowPerPage = 10
)

func init() {
	renderOpts = render.Options{
		Layout:     "main",
		Directory:  "views",
		Extensions: []string{".html"},
		Funcs: []template.FuncMap{
			{
				"add": func(x int, y int) int {
					return x + y
				},
				"isset": isSet,
				"copy": func(src, dest int) int {
					dest = src
					return dest
				},
				"isNil": func(a interface{}) int {
					b := reflect.ValueOf(a)
					switch b.Kind() {
					case reflect.Interface:
						if reflect.ValueOf(a).IsNil() {
							return 1
						} else {
							return 0
						}
					case reflect.Struct:
						return 0
					}
					return 0
				},
			},
		},
	}

	renderOpts.IsDevelopment = true
	if os.Getenv("ENV") == "production" {
		renderOpts.IsDevelopment = false
	}
	gob.Register(map[string]interface{}{})
}

func FatalnOnErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

type Response struct {
	Title   string      `json:"title"`
	BaseUrl string      `json:"base_url"`
	AppName string      `json:"app_name"`
	Error   int         `json:"error"`
	Message interface{} `json:"message"`
	Data    interface{} `json:"data"`
}

func NewResponse() *Response {
	return &Response{
		Title:   "Welcome",
		BaseUrl: "/web",
		AppName: "Gofeed",
		Error:   0,
		Message: "",
		Data:    "",
	}
}

func main() {
	configFile := "config/database.json"
	if os.Getenv("GOFEED_CONFIG") != "" {
		configFile = os.Getenv("GOFEED_CONFIG")
	}
	config := getSetting(configFile)
	connStr := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True",
		config.Username,
		config.Password,
		config.Hostname,
		config.Port,
		config.DbName,
	)

	var err error
	db, err = gorm.Open("mysql", connStr)
	if err != nil {
		log.Fatalln(err)
	}

	db.DB()
	db.DB().Ping()
	db.DB().SetMaxIdleConns(10)
	db.SingularTable(true)

	r := mux.NewRouter()
	s := http.StripPrefix("/web/", http.FileServer(http.Dir("web")))
	r.PathPrefix("/web/").Handler(s)

	r.HandleFunc("/page/{page:[0-9]+}", indexHandler).Methods("GET")
	r.HandleFunc("/", indexHandler).Methods("GET")

	// site
	site := r.PathPrefix("/site").Subrouter()
	site.HandleFunc("/view/{id:[0-9]+}", siteView)
	site.HandleFunc("/login", siteLogin).Methods("GET")
	site.HandleFunc("/login", siteLoginPost).Methods("POST")
	site.HandleFunc("/logout", siteLogout).Methods("GET")
	site.HandleFunc("/signup", siteSignup).Methods("GET")
	site.HandleFunc("/signup", siteSignupPost).Methods("POST")
	site.HandleFunc("/about", siteAbout).Methods("GET")
	site.HandleFunc("/unauthorized", siteUnauthorized).Methods("GET")

	// news
	news := r.PathPrefix("/news").Subrouter()
	news.Handle("/index", negroni.New(
		negroni.HandlerFunc(userMiddleware),
		negroni.HandlerFunc(newsIndex),
	))
	news.Handle("/suggest", negroni.New(
		negroni.HandlerFunc(userMiddleware),
		negroni.HandlerFunc(newsSuggest),
	))
	news.Handle("/admin", negroni.New(
		negroni.HandlerFunc(userMiddleware),
		negroni.HandlerFunc(newsAdmin),
	))
	news.Handle("/edit/{id:[0-9]+}", negroni.New(
		negroni.HandlerFunc(userMiddleware),
		negroni.HandlerFunc(newsEdit),
	)).Methods("GET")
	news.Handle("/update/{id:[0-9]+}", negroni.New(
		negroni.HandlerFunc(userMiddleware),
		negroni.HandlerFunc(newsUpdate),
	)).Methods("POST")
	news.Handle("/delete/{id:[0-9]+}", negroni.New(
		negroni.HandlerFunc(userMiddleware),
		negroni.HandlerFunc(newsDelete),
	))
	news.Handle("/view/{id:[0-9]+}", negroni.New(
		negroni.HandlerFunc(userMiddleware),
		negroni.HandlerFunc(newsView),
	))
	news.HandleFunc("/rss", newsRss)
	news.HandleFunc("/atom", newsAtom)

	// comment
	comment := r.PathPrefix("/comment").Subrouter()
	comment.Handle("/index", negroni.New(
		negroni.HandlerFunc(userMiddleware),
		negroni.HandlerFunc(commentIndex),
	))
	comment.Handle("/delete/{id:[0-9]+}", negroni.New(
		negroni.HandlerFunc(userMiddleware),
		negroni.HandlerFunc(commentDelete),
	))

	// user
	user := r.PathPrefix("/user").Subrouter()
	user.Handle("/index", negroni.New(
		negroni.HandlerFunc(userMiddleware),
		negroni.HandlerFunc(userIndex),
	))

	n := negroni.Classic()
	store := cookiestore.New([]byte("ohhhsooosecret"))
	n.Use(sessions.Sessions("global_session_store", store))
	n.UseHandler(r)

	port := "3000"
	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	}
	graceful.Run(":"+port, 1*time.Second, n)
}

func indexHandler(w http.ResponseWriter, req *http.Request) {
	session := sessions.GetSession(req)

	var count int
	db.Table("news").Count(&count)
	totalPage := count / rowPerPage

	params := mux.Vars(req)
	p := 1
	offset := 0
	if params["page"] != "" {
		page, _ := strconv.Atoi(params["page"])
		p = page
		if p < 0 || p > totalPage {
			p = 1
		}
		offset = (p - 1) * rowPerPage
	}

	renderOpts.Layout = "main"
	rnd = render.New(renderOpts)
	response := NewResponse()
	response.Title = "News"

	linkPage := make([]int, 0)
	news := make([]News, 0)

	data := make(map[string]interface{})
	var u User
	sUser := session.Get("user")
	if sUser != nil {
		if errUser := json.Unmarshal(sUser.([]byte), &u); errUser == nil {
			data["user"] = u
		}
	}
	if db.Table("news").Offset(offset).Limit(10).Find(&news).RecordNotFound() == false {
		for i := 1; i <= totalPage; i++ {
			linkPage = append(linkPage, i)
		}
		data["news"] = news
		data["totalPage"] = totalPage
		data["linkPage"] = linkPage
		data["currentPage"] = p
		response.Data = data
	}
	rnd.HTML(w, http.StatusOK, "site/index", response)
}

func getSetting(fileName string) DbConfig {
	val, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatalln(err)
	}

	var conf DbConfig
	err = json.Unmarshal(val, &conf)
	if err != nil {
		log.Fatalln(err)
	}
	return conf
}
