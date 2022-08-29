package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func Static(pattern, root string) {
	if _, err := os.Stat(root); os.IsNotExist(err) {
		log.Println(err.Error())
		return
	}

	http.Handle(pattern, http.FileServer(http.Dir(root)))
}

func Start(port int) {
	p := fmt.Sprintf(":%d", port)
	log.Fatal(http.ListenAndServe(p, nil))
}
