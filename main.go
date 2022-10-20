package main

import (
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(wri http.ResponseWriter, req *http.Request) {
		wri.Write([]byte(`
			<html>
			   <head>
				 <title>Chat</title>
			   </head>
			   <body>
				 Let's chat!
			   </body>
			</html>
		`))
	})
	//Inicia o web service
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
