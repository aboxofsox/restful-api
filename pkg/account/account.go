package account

import (
	"database/sql"
	"fmt"
	"log"
)

type Account struct {
	Id       int    `json:"id,omitempty"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

func New(username, password string) *Account {
	return &Account{
		Username: username,
		Password: password,
	}
}

func (acc *Account) Save(db *sql.DB) error {
	es := fmt.Sprintf("insert into accounts (username, password) values ('%s', '%s');", acc.Username, acc.Password)

	res, err := db.Exec(es)
	if err != nil {
		return fmt.Errorf("account: sql exec error: %v", err)
	}

	ra, _ := res.RowsAffected()

	log.Printf("%d rows affected", ra)
	return nil
}
