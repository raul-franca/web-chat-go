package main

import (
	"flag"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"sync"
)

type TemplateHandler struct {
	once     sync.Once
	fileName string
	templ    *template.Template
}

// ServeHTTP handles the HTTP request.
func (t *TemplateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.fileName)))
	})
	t.templ.Execute(w, r)
}

func main() {

	var addr = flag.String("addr", ":8080", "The addr of the  application.")
	flag.Parse() // parse the flags

	r := newRoom()
	http.Handle("/", &TemplateHandler{fileName: "chat.html"})
	http.Handle("/room", r)
	go r.run()
	//Inicia o web service
	log.Println("Starting web server on", *addr)
	if err := http.ListenAndServe(*addr, nil)
		err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
