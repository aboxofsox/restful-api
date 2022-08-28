package server

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"rest/pkg/account"
	"rest/pkg/auth"
)

// Ok is a simple GET test route
func Ok(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	log.Println("connected:", r.RemoteAddr)

	if _, err := w.Write([]byte("hello world")); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// NewAccount creates a new account.Account given form data
func NewAccount(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		if r.FormValue("username") == "" || r.FormValue("password") == "" {
			http.Error(w, "no username or password", http.StatusNotAcceptable)
			return
		}

		username := r.FormValue("username")
		password := r.FormValue("password")

		acc := account.New(username, auth.HashAndSalt(password))

		if err := acc.Save(db); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// Login handles user authentication
func Login(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var hash string
		username := r.FormValue("username")
		password := r.FormValue("password")

		if username == "" || password == "" {
			http.Error(w, "no username or password", http.StatusNotAcceptable)
			return
		}

		qs := fmt.Sprintf("select password from accounts where username = '%s';", username)

		row := db.QueryRow(qs)
		if errors.Is(row.Err(), sql.ErrNoRows) {
			http.Error(w, "username found", http.StatusNotFound)
			return
		}

		if err := row.Scan(&hash); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if auth.Compare(hash, password) {
			if _, err := w.Write([]byte("authenticated")); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

		} else {
			http.Error(w, "not authenticated", http.StatusForbidden)
			return
		}

	}
}

// GetAllAccounts sends every account in the database as JSON
func GetAllAccounts(db *sql.DB) http.HandlerFunc {
	var acc account.Account
	accounts := make([]account.Account, 10)

	qs := "select * from accounts;"

	rows, err := db.Query(qs)
	if err != nil {
		return nil
	}

	for rows.Next() {
		if err := rows.Scan(&acc.Id, &acc.Username, &acc.Password); err != nil {
			log.Println(err.Error())
			return nil
		}

		accounts = append(accounts, acc)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		b, err := json.MarshalIndent(accounts, "", " ")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if _, err := w.Write(b); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
