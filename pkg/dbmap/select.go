package dbmap

import (
	"context"
	"fmt"
	"log"
	"reflect"

	"github.com/jmoiron/sqlx"
)

func (d DBMap) Select(ctx context.Context, data interface{}, conditions string) error {
	var stmt *sqlx.Stmt
	var rows *sqlx.Rows
	var err error
	var str string
	var colnames string
	var join string
	var order string
	var query string

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

	q := d.query[str]
	if q != nil {
		colnames = q[0]
		join = q[1]
		order = q[2]
	}
	if order != "" {
		query = fmt.Sprintf("SELECT %s FROM %s WHERE %s ORDER BY %s", colnames, join, conditions, order)
	} else {
		query = fmt.Sprintf("SELECT %s FROM %s WHERE %s", colnames, join, conditions)
	}
	log.Println(query)

	if ctx != nil {
		stmt, err = d.db.PreparexContext(ctx, query)
	} else {
		stmt, err = d.db.Preparex(query)
	}

	if err != nil {
		return err
	}

	switch returnType {
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

		data = arr
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
			err := rows.StructScan(res)
			if err == nil {
				ar.Set(reflect.Append(ar, reflect.ValueOf(res)))
			}
			res = reflect.New(reflect.TypeOf(res).Elem()).Interface()
		}

		a, total := d.ValueToMapString(ar.Interface())
		d.MapRecursive(data, a, total)

		data = arr
	}
	return nil
}
