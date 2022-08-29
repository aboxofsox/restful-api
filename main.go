package main

import (
	"os"
	"rest/pkg/server"
	"strconv"
)

func main() {
	args := os.Args[1:]
	p, _ := strconv.Atoi(args[0])
	server.Static("/", "public")
	server.Start(p)
}
