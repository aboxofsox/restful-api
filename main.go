package main

import (
	"log"
	"rest/pkg/server"
)

const DBConnection string = ""

func main() {
	s := server.New(3000)

	if err := s.ConnectDB("postgres", DBConnection); err != nil {
		log.Fatal(err.Error())
	}

	s.AddRoute("/ok", server.Ok)
	s.AddRoute("/new-account", server.NewAccount(s.DB()))
	s.AddRoute("/login", server.Login(s.DB()))
	s.AddRoute("/get-all", server.GetAllAccounts(s.DB()))

	if err := s.Start(); err != nil {
		log.Fatal(err.Error())
	}
}
