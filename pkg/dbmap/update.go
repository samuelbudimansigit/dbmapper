package dbmap

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"reflect"
)

func (d DBMap) Update(ctx context.Context, data interface{}) error {
	tableName := reflect.TypeOf(data).Elem()
	var m *Model

	m = d.models[fmt.Sprintf("%v", tableName)]

	if m == nil {
		return errors.New("table not found")
	}

	var updatedId []interface{}
	var bufferColumn bytes.Buffer
	var bufferCondition bytes.Buffer

	var val []interface{}

	t := reflect.ValueOf(data).Elem()

	for i := 0; i < t.NumField(); i++ {
		varName := t.Type().Field(i).Tag.Get("db")
		varValue := t.Field(i).Interface()
		if contains(varName, m.primaryKey) {
			updatedId = append(updatedId, varValue)
		} else if varName != "" {
			val = append(val, varValue)
		}
	}

	flag := 1

	for _, col := range m.columns {
		if !contains(col.name, m.primaryKey) {
			bufferColumn.WriteString(fmt.Sprintf("%s = $%d,", col.name, flag))
			flag++
		}
	}

	for _, pk := range m.primaryKey {
		bufferCondition.WriteString(fmt.Sprintf("%v = $%d AND", pk, flag))
		flag++
	}

	for _, v := range updatedId {
		val = append(val, v)
	}

	columns := bufferColumn.String()
	columns = columns[:len(columns)-1]

	conditions := bufferCondition.String()
	conditions = conditions[:len(conditions)-4]

	query := fmt.Sprintf("UPDATE %s SET %s WHERE %s", m.table, columns, conditions)

	if ctx != nil {
		_ = d.db.MustExecContext(ctx, query, val...)
	} else {
		_ = d.db.MustExec(query, val...)
	}

	return nil
}
