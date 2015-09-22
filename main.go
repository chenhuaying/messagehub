package main

import (
	_ "fmt"
	"html/template"
	"log"
	"net/http"
	_ "strconv"
)

var (
	homeTempl = template.Must(template.ParseFiles("page.html"))
)

func handler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "Not found", 404)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}

	var v = struct {
		Host string
		Data string
	}{
		r.Host,
		"Testing WebSocket",
	}
	if err := homeTempl.ExecuteTemplate(w, "page", v); err != nil {
		log.Println("ExecuteTemplate, error: ", err)
	}
}

func main() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/ws", handlerws)
	go tracker.run()
	go http.ListenAndServe(":8080", nil)
	http.ListenAndServeTLS(":8081", "server.crt", "server.key", nil)
}
