package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/feeds"
	// "github.com/unrolled/render"
)

func newsIndex(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
}

func newsSuggest(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
}

func newsAdmin(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
}

func newsEdit(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
}

func newsUpdate(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
}

func newsDelete(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
}

func newsView(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
}

func newsRss(w http.ResponseWriter, req *http.Request) {
	now := time.Now()
	feed := &feeds.Feed{
		Title:       "GoFeed",
		Link:        &feeds.Link{Href: "http://localhost:3000"},
		Description: "News all about Go Programming",
		Author:      &feeds.Author{Name: "Agus Susilo", Email: "smartgdi@gmail.com"},
		Created:     now,
	}

	feedItem := []*feeds.Item{}
	news := make([]News, 0)
	if db.Table("news").Limit(50).Where("status= ?", 1).Order("created_at desc").
		Find(&news).RecordNotFound() == false {
		for _, val := range news {
			id := strconv.FormatInt(val.ID, 10)
			feedItem = append(feedItem, &feeds.Item{
				Title:       val.Title,
				Link:        &feeds.Link{Href: "http://localhost:3000/news/" + id},
				Description: val.Text,
				Created:     val.CreatedAt,
			})
		}
	}
	feed.Items = feedItem
	atom, err := feed.ToRss()
	if err != nil {
		log.Println(err)
	}

	fmt.Fprint(w, atom)
}

func newsAtom(w http.ResponseWriter, req *http.Request) {
	now := time.Now()
	feed := &feeds.Feed{
		Title:       "GoFeed",
		Link:        &feeds.Link{Href: "http://localhost:3000"},
		Description: "News all about Go Programming",
		Author:      &feeds.Author{Name: "Agus Susilo", Email: "smartgdi@gmail.com"},
		Created:     now,
	}

	feedItem := []*feeds.Item{}
	news := make([]News, 0)
	if db.Table("news").Limit(50).Where("status= ?", 1).Order("created_at desc").
		Find(&news).RecordNotFound() == false {
		for _, val := range news {
			id := strconv.FormatInt(val.ID, 10)
			feedItem = append(feedItem, &feeds.Item{
				Title:       val.Title,
				Link:        &feeds.Link{Href: "http://localhost:3000/news/" + id},
				Description: val.Text,
				Created:     val.CreatedAt,
			})
		}
	}
	feed.Items = feedItem
	atom, err := feed.ToAtom()
	if err != nil {
		log.Println(err)
	}

	fmt.Fprint(w, atom)
}
