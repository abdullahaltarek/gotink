package main

import (
	"github.com/cyantarek/gotink"
	"net/http"
)

type Hello struct {
	Msg string 	`json:"msg"`
}

func main() {
	app := tink.New()
	app.Port = "8021"
	app.Log = true
	app.Json("GET", "/aaa", HttpHandler)
	app.Json("GET", "/bbb", JsonHandler)
	app.Run()
}

func HttpHandler(w http.ResponseWriter, r *http.Request) {
	tink.RenderHtml(w, "index.html", nil)
}


func JsonHandler(w http.ResponseWriter, r *http.Request) {
	msg := Hello{Msg:"Hello"}
	tink.RenderJson(w, msg)
}