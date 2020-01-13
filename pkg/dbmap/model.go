package dbmap

import (
	"github.com/jmoiron/sqlx"
)

type DBMap struct {
	db        *sqlx.DB
	models    map[string]*Model
	relations map[string][]*Relation
	Flats     map[string]interface{}
	query     map[string][]string
}

type Model struct {
	table      string
	columns    []Column
	primaryKey []string
	autoKey    bool
	data       interface{}
}

type Column struct {
	name     string
	dataType string
}

type Relation struct {
	from         interface{}
	to           interface{}
	pK           interface{}
	fK           interface{}
	by           interface{}
	relationType int
}
