package main

import "database/sql"

type State struct {
	db *sql.DB
}

func (s *State) GetDB() *sql.DB {
	return s.db
}

func NewState() *State {
	return &State{}
}
