package dbmap

import (
	"fmt"
	"reflect"
)

func (d DBMap) AddTable(tableName string, primaryKey []string, autoKey bool, data interface{}) error {
	var columns []Column
	t := reflect.ValueOf(data)

	for i := 0; i < t.NumField(); i++ {
		varTag := t.Type().Field(i).Tag.Get("db")
		if varTag != "" {
			column := Column{
				name: varTag,
			}
			columns = append(columns, column)
		}

	}

	model := &Model{
		table:      tableName,
		columns:    columns,
		primaryKey: primaryKey,
		autoKey:    autoKey,
		data:       data,
	}

	tableName = fmt.Sprintf("%v", reflect.TypeOf(data))
	d.models[tableName] = model

	return nil
}
