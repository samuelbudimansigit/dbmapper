package dbmap

import (
	"bytes"
	"fmt"
	"reflect"
)

func (d DBMap) CreateQuery(data interface{}) (interface{}, string, string) {
	fs := []reflect.StructField{}

	var col bytes.Buffer

	t := reflect.ValueOf(data)
	m := d.models[fmt.Sprintf("%v", reflect.TypeOf(data))]

	for i := 0; i < t.NumField(); i++ {
		if t.Type().Field(i).Tag.Get("db") != "" {
			sf := reflect.StructField{
				Name: fmt.Sprintf("%s%s", upperFirst(m.table), t.Type().Field(i).Name),
				Type: t.Type().Field(i).Type,
				Tag:  reflect.StructTag(fmt.Sprintf(`db:"%s.%s"`, m.table, t.Type().Field(i).Tag.Get("db"))),
			}
			col.WriteString(sf.Tag.Get("db") + " as " + `"` + sf.Tag.Get("db") + `"` + ",")
			fs = append(fs, sf)
		}
	}

	v := reflect.New(reflect.StructOf(fs)).Elem()
	s := v.Addr().Interface()

	colRes := col.String()
	colRes = colRes[:len(colRes)-1]

	return s, colRes, m.table
}
