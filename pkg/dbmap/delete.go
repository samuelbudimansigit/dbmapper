package dbmap

import (
	"bytes"
	"context"
	"fmt"
	"reflect"
	"strings"
)

func (d DBMap) Delete(ctx context.Context, data interface{}) error {
	var deletedId []interface{}
	var bufferCondition bytes.Buffer

	tableName := fmt.Sprintf("%v", reflect.TypeOf(data).Elem())
	tableName = strings.Replace(tableName, "[]", "", 1)
	m := d.models[tableName]

	t := reflect.ValueOf(data).Elem()

	for i := 0; i < t.NumField(); i++ {
		varName := t.Type().Field(i).Tag.Get("db")
		if contains(varName, m.primaryKey) {
			deletedId = append(deletedId, t.Field(i).Interface())
		}
	}

	for i, pk := range m.primaryKey {
		bufferCondition.WriteString(fmt.Sprintf("%s = $%d AND ", pk, (i + 1)))
	}

	conditions := bufferCondition.String()
	conditions = conditions[:len(conditions)-4]

	query := fmt.Sprintf("DELETE FROM %s WHERE %s", m.table, conditions)

	if ctx != nil {
		_ = d.db.MustExecContext(ctx, query, deletedId...)
	} else {
		_ = d.db.MustExec(query, deletedId...)
	}

	return nil
}
