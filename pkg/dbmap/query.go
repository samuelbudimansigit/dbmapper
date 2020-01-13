package dbmap

import (
	"fmt"
	"reflect"

	"github.com/jmoiron/sqlx"
	"golang.org/x/net/context"
)

type Flat struct {
	PersonId      int    `db:"person.id"`
	PersonName    string `db:"person.name"`
	PersonJobId   int    `db:"person.job_id"`
	PacarId       int    `db:"pacar.id"`
	PacarName     string `db:"pacar.name"`
	PacarPersonId int    `db:"pacar.person_id"`
	ShoesId       int    `db:"shoes.id"`
	ShoesName     string `db:"shoes.name"`
	ShoesPacarId  int    `db:"shoes.pacar_id"`
	JobId         int    `db:"job.id"`
	JobName       string `db:"job.name"`
}

func (d DBMap) Query(ctx context.Context, data interface{}, query string) error {
	var stmt *sqlx.Stmt
	var rows *sqlx.Rows
	var str string
	var err error

	tableName := reflect.TypeOf(data).Elem()
	returnType := 1
	var m *Model

	m = d.models[fmt.Sprintf("%v", tableName)]

	if m == nil {
		str = fmt.Sprintf("%v", reflect.TypeOf(data).Elem().Elem())
		m = d.models[fmt.Sprintf("%v", reflect.TypeOf(data).Elem().Elem())]
		returnType = 2
	}

	if m == nil {
		str = fmt.Sprintf("%v", reflect.TypeOf(data).Elem().Elem().Elem())
		m = d.models[fmt.Sprintf("%v", reflect.TypeOf(data).Elem().Elem().Elem())]
	}

	if ctx != nil {
		stmt, err = d.db.PreparexContext(ctx, query)
	} else {
		stmt, err = d.db.Preparex(query)
	}

	if err != nil {
		return err
	}

	switch returnType {
	case type_pointer:
		res := d.Flats[str]
		arr := reflect.TypeOf(res)
		ar := reflect.New(reflect.SliceOf(arr)).Elem()

		if ctx != nil {
			rows, err = stmt.QueryxContext(ctx)
		} else {
			rows, err = stmt.Queryx()
		}

		if err != nil {
			return err
		}

		for rows.Next() {
			if err = rows.StructScan(res); err == nil {
				ar.Set(reflect.Append(ar, reflect.ValueOf(res)))
			}
			res = reflect.New(reflect.TypeOf(res).Elem()).Interface()
		}

		a, total := d.ValueToMapString(ar.Interface())
		d.MapRecursive(data, a, total)

	case type_slice:
		res := d.Flats[str]
		arr := reflect.TypeOf(res)
		ar := reflect.New(reflect.SliceOf(arr)).Elem()

		if ctx != nil {
			rows, err = stmt.QueryxContext(ctx)
		} else {
			rows, err = stmt.Queryx()
		}

		if err != nil {
			return err
		}

		for rows.Next() {
			err := rows.StructScan(res)
			if err == nil {
				ar.Set(reflect.Append(ar, reflect.ValueOf(res)))
			}
			res = reflect.New(reflect.TypeOf(res).Elem()).Interface()
		}

		a, total := d.ValueToMapString(ar.Interface())
		d.MapRecursive(data, a, total)

	}
	return nil
}
