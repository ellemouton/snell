package main

import (
	"database/sql"
	"fmt"

	"github.com/ellemouton/snell/db"
	"github.com/ellemouton/snell/lnd"
	"github.com/ellemouton/snell/macaroon"
)

type State struct {
	db        *sql.DB
	lndClient lnd.Client
	macClient macaroon.Client
}

func NewState() (*State, error) {
	s := new(State)

	db, err := db.Connect()
	if err != nil {
		return nil, fmt.Errorf("problem connecting to db: %s", err)
	}
	s.db = db

	mc, err := macaroon.New()
	if err != nil {
		return nil, fmt.Errorf("problem creating macaroon client: %s", err)
	}
	s.macClient = mc

	lc, err := lnd.New()
	if err != nil {
		return nil, fmt.Errorf("problem creating lnd client: %s", err)
	}
	s.lndClient = lc

	return s, nil
}

func (s *State) GetDB() *sql.DB {
	return s.db
}

func (s *State) GetMacaroonClient() macaroon.Client {
	return s.macClient
}

func (s *State) GetLndClient() lnd.Client {
	return s.lndClient
}

func (s *State) cleanup() {
	if err := s.db.Close(); err != nil {
		fmt.Errorf("error closing db: %v", err)
	}

	if err := s.macClient.Close(); err != nil {
		fmt.Errorf("error closing db: %v", err)
	}

	if err := s.lndClient.Close(); err != nil {
		fmt.Errorf("error closing db: %v", err)
	}

}
