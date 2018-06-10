package tink

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"time"
)

type (
	App struct {
		Mux  *http.ServeMux
		Log  bool
		Port string
		TemplateDir string
	}
)

const (
	defaultPort = "8000"
)

func New() *App {
	app := App{Mux: http.NewServeMux()}

	return &app
}

func (a *App) Run() {
	log.Println("Starting Http Server...")

	if a.Port != "" {
		log.Println("Running on http://127.0.0.1:" + a.Port)
	} else {
		log.Println("Running on http://127.0.0.1:" + defaultPort)
	}

	if a.Log && a.Port != "" {
		log.Println("Logger enabled")
		http.ListenAndServe(":"+a.Port, logRequest(a.Mux))
	} else if a.Log && a.Port == ""  {
		http.ListenAndServe(":"+defaultPort, logRequest(a.Mux))
	} else if !a.Log && a.Port == ""  {
		http.ListenAndServe(":"+ defaultPort, a.Mux)
	} else {
		http.ListenAndServe(":" + a.Port, logRequest(a.Mux))
	}

}

func (a *App) Http(method string, path string, handlerFunc http.HandlerFunc) {
	switch method {
	case "GET":
		handlerFunc = get(handlerFunc)
		handlerFunc = httpResponse(handlerFunc)
		a.Mux.HandleFunc(path, handlerFunc)

	case "POST":
		handlerFunc = post(handlerFunc)
		handlerFunc = httpResponse(handlerFunc)
		a.Mux.HandleFunc(path, handlerFunc)

	case "PUT":
		handlerFunc = put(handlerFunc)
		handlerFunc = httpResponse(handlerFunc)
		a.Mux.HandleFunc(path, handlerFunc)

	case "DELETE":
		handlerFunc = delete_(handlerFunc)
		handlerFunc = httpResponse(handlerFunc)
		a.Mux.HandleFunc(path, handlerFunc)
	}
}

func (a *App) Json(method string, path string, handlerFunc http.HandlerFunc) {
	switch method {
	case "GET":
		handlerFunc = get(handlerFunc)
		handlerFunc = jsonResponse(handlerFunc)
		a.Mux.HandleFunc(path, handlerFunc)

	case "POST":
		handlerFunc = post(handlerFunc)
		handlerFunc = jsonResponse(handlerFunc)
		a.Mux.HandleFunc(path, handlerFunc)

	case "PUT":
		handlerFunc = put(handlerFunc)
		handlerFunc = jsonResponse(handlerFunc)
		a.Mux.HandleFunc(path, handlerFunc)

	case "DELETE":
		handlerFunc = delete_(handlerFunc)
		handlerFunc = jsonResponse(handlerFunc)
		a.Mux.HandleFunc(path, handlerFunc)
	}
}

func get(f func(w http.ResponseWriter, r *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			w.WriteHeader(405)
			fmt.Fprint(w, "Method not allowed")
		}
		f(w, r)
	}
}

func post(f func(w http.ResponseWriter, r *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			w.WriteHeader(405)
			fmt.Fprint(w, "Method not allowed")
		}
		f(w, r)
	}
}

func put(f func(w http.ResponseWriter, r *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			w.WriteHeader(405)
			fmt.Fprint(w, "Method not allowed")
		}
		f(w, r)
	}
}

func delete_(f func(w http.ResponseWriter, r *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			w.WriteHeader(405)
			fmt.Fprint(w, "Method not allowed")
		}
		f(w, r)
	}
}

//MIDDLE-WARES

func logRequest(f http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Println("---> Request Method: "+r.Method, "Protocol: "+r.Proto, "Path: "+r.URL.Path)
		f.ServeHTTP(w, r)
		diff := time.Since(start)
		log.Print("<--- Response Time: ", diff.Seconds())


	})
}

func httpResponse(f func(w http.ResponseWriter, r *http.Request)) func(http.ResponseWriter, *http.Request) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		f(w, r)
	})
}

func jsonResponse(f func(w http.ResponseWriter, r *http.Request)) func(http.ResponseWriter, *http.Request) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		f(w, r)
	})
}

func RenderHtml(w io.Writer, templ string, data interface{}) {
	tpl := template.Must(template.ParseFiles("./template/" + templ))
	tpl.Execute(w, data)
}

func RenderJson(w io.Writer, data interface{}) {
	json.NewEncoder(w).Encode(data)
}
