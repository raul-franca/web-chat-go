package main

import (
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
	t.templ.Execute(w, nil)
}

func main() {
	//root
	r := newRoom()
	http.Handle("/", &TemplateHandler{fileName: "chat.html"})
	http.Handle("/room", r)
	go r.run()
	//Inicia o web service
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
