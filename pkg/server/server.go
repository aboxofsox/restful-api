package server

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
)

type Server struct {
	port   string
	mux    *http.ServeMux
	routes map[string]http.HandlerFunc
	db     *sql.DB
}

func New(port int) *Server {
	return &Server{
		port:   fmt.Sprintf(":%d", port),
		mux:    http.NewServeMux(),
		routes: map[string]http.HandlerFunc{},
	}
}

func (s *Server) ConnectDB(driver, dbUrl string) error {
	db, err := sql.Open(driver, dbUrl)
	if err != nil {
		return fmt.Errorf("server: database connection error: %v", err)
	}

	if err := db.Ping(); err != nil {
		return fmt.Errorf("server: database ping error: %v", err)
	}

	return nil
}

func (s *Server) DB() *sql.DB { return s.db }

func (s *Server) AddRoute(pattern string, handler http.HandlerFunc) {
	if _, exists := s.routes[pattern]; exists {
		fmt.Printf("server: route exists error: %s", pattern)
	}

	s.routes[pattern] = handler
}

func (s *Server) Static(pattern, root string) {
	if _, exists := s.routes[pattern]; exists {
		fmt.Printf("server: route pattern exists: %s\n", pattern)
	}

	if _, err := os.Stat(root); os.IsNotExist(err) {
		fmt.Printf("server: root not exist error: %v\n", err)
	}

	fs := http.FileServer(http.Dir(root))
	http.Handle(pattern, fs)
}

func (s *Server) Start() error {
	if len(s.routes) != 0 {
		for p, fn := range s.routes {
			s.mux.HandleFunc(p, fn)
		}
	}

	log.Println("server starting", s.port)

	if err := http.ListenAndServe(s.port, s.mux); err != nil {
		return fmt.Errorf("server: start error: %v", err)
	}

	return nil
}