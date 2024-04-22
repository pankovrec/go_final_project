package main

import (
	"net/http"
	"db"
	"repeat"
//	"github.com/go-chi/chi/v5"
)

const (
	webDir string = "./web"
)

func main() {
	db.StartDb()
//	r := chi.NewRouter()
	//r.Mount("/", http.FileServer(http.Dir(webDir)))
	//r.Get("/api/nextdate", repeat.ExtractParamsDate)

	startWeb()
	
}

func startWeb() {
	http.Handle("/", http.FileServer(http.Dir(webDir)))
	http.Handle("/api/nextdate", http.HandlerFunc(repeat.ExtractParamsDate))

//http.Handle("/", http.FileServer(http.Dir(webDir)))

err := http.ListenAndServe(":7540", nil)
if err != nil {
	panic(err)
}
}
