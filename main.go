package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"rest/pkg/server"
	"strings"
)

func rando() string {
	bt := make([]byte, 256)
	_, err := rand.Read(bt)
	if err != nil {
		return ""
	}

	return hex.EncodeToString(bt)
}

func getRando(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(rando()))
}

func userId(w http.ResponseWriter, r *http.Request) {
	id := server.Params(r, "userid")
	fmt.Println(id)
}

func message(w http.ResponseWriter, r *http.Request) {
	msg := server.Params(r, "message")
	fmt.Println(msg)
}

func reverse(w http.ResponseWriter, r *http.Request) {
	arg := server.Params(r, "arg")

	argSplit := strings.Split(arg, "")

	s := make([]string, len(argSplit))

	for i := len(s) - 1; i >= 0; i-- {
		s = append(s, argSplit[i])
	}

	fmt.Println(strings.Join(s, ""))
}

func main() {
	router := &server.Router{}
	router.Route("GET", "/api/rando", getRando)
	router.Route("GET", `/api/users/(?P<userid>\d+)`, userId)
	router.Route("GET", `/api/message/(?P<message>\w+)`, message)
	router.Route("GET", `/api/reverse/(?P<arg>\w+)`, reverse)
	server.Start(3000, router)
}
