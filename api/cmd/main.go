package main

import (
	"log"
	"maga-auctions/api"
	"maga-auctions/utils"
	"net/http"
	"time"
)

func main() {
	port := ":" + utils.EnvVars.API.Port

	if port == ":" {
		log.Fatal("PORT must be set")
	}

	s := &http.Server{
		Addr:           port,
		Handler:        api.Config(),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	log.Panic("error at listen and serve", s.ListenAndServe())
}
