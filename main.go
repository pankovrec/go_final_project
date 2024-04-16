package main

import (
	"net/http"
	"db"
)

const (
	webDir string = "./web"
)

func main() {
	db.StartDb()

	startWeb()

}

func startWeb() {
http.Handle("/", http.FileServer(http.Dir(webDir)))

err := http.ListenAndServe(":7540", nil)
if err != nil {
	panic(err)
}
}
