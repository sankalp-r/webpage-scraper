package router

import (
	"github.com/gorilla/mux"
	"h24/pkg/handlers"
	"net/http"
	"time"
)

func Getrouter() http.Handler{
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", handlers.HomePage)
	myRouter.HandleFunc("/crawl",handlers.Crawl)
	muxWithMiddlewares := http.TimeoutHandler(myRouter, time.Second*10, "Timeout!")
	return muxWithMiddlewares
}
