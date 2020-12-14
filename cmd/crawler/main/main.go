package main

import (
	"fmt"
	router2 "h24/pkg/router"
	"log"
	"net/http"
	"time"
)


func main(){
	fmt.Println("server started...")

	router := router2.Getrouter()

	srv := &http.Server{
		ReadTimeout: 5 * time.Second,
		WriteTimeout: 10 * time.Second,
		Addr: ":10000",
		Handler: router,
	}
	log.Fatal(srv.ListenAndServe())
}




