package dbmap

import (
	"strings"
	"unicode"

	"github.com/jmoiron/sqlx"
	"github.com/jmoiron/sqlx/reflectx"
)

func CreateDBMap(db *sqlx.DB) *DBMap {
	return &DBMap{
		db:        db,
		models:    map[string]*Model{},
		relations: map[string][]*Relation{},
		Flats:     make(map[string]interface{}),
		query:     make(map[string][]string),
	}
}

func contains(word string, words []string) bool {
	for _, w := range words {
		if w == word {
			return true
		}
	}
	return false
}

func upperFirst(str string) string {
	for i, v := range str {
		return string(unicode.ToUpper(v)) + str[i+1:]
	}
	return ""
}

func (d DBMap) AddFlats(name string, flat interface{}, q1 string, q2 string, q3 string) {
	d.Flats[name] = flat
	d.query[name] = append(d.query[name], q1)
	d.query[name] = append(d.query[name], q2)
	d.query[name] = append(d.query[name], q3)
}

func (d DBMap) GetFlat(name string) interface{} {
	return d.Flats[name]
}

func (d DBMap) ChangeTag(tag string) {
	d.db.Mapper = reflectx.NewMapperFunc(tag, strings.ToLower)
}

func (d DBMap) Init() {
	for tag, rel := range d.relations {
		for _, r := range rel {
			f, col, join, o := d.CreateStruct(r.from)
			d.AddFlats(tag, f, col, join, o)
			break
		}
	}

	for tag, mod := range d.models {
		if d.relations[tag] == nil {
			f, col, join := d.CreateQuery(mod.data)
			d.AddFlats(tag, f, col, join, "")
		}
	}
}
