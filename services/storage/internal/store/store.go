package store

import (
	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"time"
)

const defaultTimeout = time.Duration(time.Second * 15)

type Store struct {
	db      *sqlx.DB
	builder squirrel.StatementBuilderType
}

func NewStore(db *sqlx.DB) *Store {
	return &Store{
		db:      db,
		builder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}
