package db

import (
	"errors"
	"fmt"

	r "gopkg.in/rethinkdb/rethinkdb-go.v6"
)

var TokenFail = errors.New("failed to get a token")

type Token struct {
	Value string
}

var TokenTable = r.DB("config").TableCreate("token")

func (d *DB) Token() (*Token, error) {
	_f := "(*DB).Token"
	Log.Debug(_f, "enter")

	d.Once(ConfigDB)
	d.Once(TokenTable)

	toks := make([]*Token, 0)
	q := r.DB("config").Table("token")
	err := q.ReadAll(&toks, d)
	if err != nil {
		err = fmt.Errorf("readall &toks: %w", err)
		Log.Error(_f, err)
		return nil, err
	}

	if len(toks) == 0 {
		return nil, TokenFail
	}

	return toks[0], nil
}
