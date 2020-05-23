package main

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/ellemouton/snell/db"
)

type State struct {
	db *sql.DB
}

func NewState() (*State, error) {
	s := &State{}

	db, err := db.Connect()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("problem connecting to db: %s", err.Error()))
	}
	s.db = db

	return s, nil
}

func (s *State) GetDB() *sql.DB {
	return s.db
}
