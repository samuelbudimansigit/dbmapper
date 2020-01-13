package dbmap

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"reflect"
)

func (d DBMap) Insert(ctx context.Context, data interface{}) error {
	tableName := reflect.TypeOf(data).Elem()
	returnType := 1
	var m *Model

	m = d.models[fmt.Sprintf("%v", tableName)]

	if m == nil {
		m = d.models[fmt.Sprintf("%v", reflect.TypeOf(data).Elem().Elem())]
		returnType = 2
	}

	if m == nil {
		return errors.New("table not found")
	}

	var bufferValues bytes.Buffer
	var bufferColumn bytes.Buffer
	var val []interface{}

	switch returnType {

	case type_slice:
		var bufferBulk bytes.Buffer
		st := reflect.ValueOf(data).Elem()

		for _, col := range m.columns {
			if m.autoKey == true {
				if !contains(col.name, m.primaryKey) {
					bufferColumn.WriteString(col.name + ",")
				}
			} else {
				bufferColumn.WriteString(col.name + ",")
			}
		}

		columns := bufferColumn.String()
		columns = columns[:len(columns)-1]
		flag := 1

		for j := 0; j < st.Len(); j++ {
			x := st.Index(j)
			if m.autoKey == false {
				for i := 0; i < x.NumField(); i++ {
					if x.Type().Field(i).Tag.Get("db") != "" {
						varValue := x.Field(i).Interface()
						bufferValues.WriteString(fmt.Sprintf("%s%d,", "$", flag))
						val = append(val, varValue)
						flag++
					}
				}
			} else {
				for i := 0; i < x.NumField(); i++ {
					varName := x.Type().Field(i).Tag.Get("db")
					if varName != "" {
						if !contains(string(varName), m.primaryKey) {
							varValue := x.Field(i).Interface()
							bufferValues.WriteString(fmt.Sprintf("%s%d,", "$", flag))
							val = append(val, varValue)
							flag++
						}
					}
				}
			}

			values := bufferValues.String()
			values = values[:len(values)-1]

			bufferBulk.WriteString(fmt.Sprintf("(%s),", values))
			bufferValues.Reset()
		}

		bulk := bufferBulk.String()
		bulk = bulk[:len(bulk)-1]

		query := fmt.Sprintf("INSERT INTO %s(%s) VALUES %s", m.table, columns, bulk)

		if ctx != nil {
			_ = d.db.MustExecContext(ctx, query, val...)
		} else {
			_ = d.db.MustExec(query, val...)
		}

	case type_pointer:
		t := reflect.ValueOf(data).Elem()

		flag := 1
		if m.autoKey == false {
			for i := 0; i < t.NumField(); i++ {
				if t.Type().Field(i).Tag.Get("db") != "" {
					varValue := t.Field(i).Interface()
					bufferValues.WriteString(fmt.Sprintf("%s%d,", "$", flag))
					val = append(val, varValue)
					flag++
				}
			}
		} else {
			for i := 0; i < t.NumField(); i++ {
				varName := t.Type().Field(i).Tag.Get("db")
				if varName != "" {
					if !contains(string(varName), m.primaryKey) {
						varValue := t.Field(i).Interface()
						bufferValues.WriteString(fmt.Sprintf("%s%d,", "$", flag))
						val = append(val, varValue)
						flag++
					}
				}
			}
		}

		values := bufferValues.String()
		values = values[:len(values)-1]

		for _, col := range m.columns {
			if m.autoKey == true {
				if !contains(col.name, m.primaryKey) {
					bufferColumn.WriteString(col.name + ",")
				}
			} else {
				bufferColumn.WriteString(col.name + ",")
			}

		}
		columns := bufferColumn.String()
		columns = columns[:len(columns)-1]

		query := fmt.Sprintf("INSERT INTO %s(%s) VALUES (%s)", m.table, columns, values)

		if ctx != nil {
			_ = d.db.MustExecContext(ctx, query, val...)
		} else {
			_ = d.db.MustExec(query, val...)
		}
	}

	return nil
}
