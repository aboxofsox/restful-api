package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

// Static creates an http.FileServer for static content
func Static(pattern, root string) {
	if _, err := os.Stat(root); os.IsNotExist(err) {
		log.Println(err.Error())
		return
	}

	http.Handle(pattern, http.FileServer(http.Dir(root)))
}

// Start starts an HTTP server on the specified port
func Start(port int, handler http.Handler) {
	p := fmt.Sprintf(":%d", port)

	log.Println("listening", p)

	log.Fatal(http.ListenAndServe(p, handler))
}
