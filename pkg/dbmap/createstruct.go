package dbmap

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"
)

func (d DBMap) CreateStruct(data interface{}) (interface{}, string, string, string) {
	fs := []reflect.StructField{}

	var col bytes.Buffer
	var join bytes.Buffer
	var order bytes.Buffer

	t := reflect.ValueOf(data)
	m := d.models[fmt.Sprintf("%v", reflect.TypeOf(data))]
	r := d.relations[fmt.Sprintf("%v", reflect.TypeOf(data))]

	join.WriteString(m.table)

	// if r != nil {
	// 	for _, rr := range r {
	// 		m2 := d.models[fmt.Sprintf("%v", reflect.TypeOf(rr.to))]
	// 		join.WriteString(" JOIN " + m2.table + " ON " + fmt.Sprintf("%s.%s", m.table, rr.pK) + "=" + fmt.Sprintf("%s.%s ", m2.table, rr.fK))
	// 	}
	// }

	for i := 0; i < t.NumField(); i++ {
		if t.Type().Field(i).Tag.Get("join") != "" {
			childName := fmt.Sprintf("%v", reflect.TypeOf(t.Field(i).Interface()))
			childName = strings.Replace(childName, "[]", "", 1)
			childName = strings.Replace(childName, "*", "", 1)
			for _, rr := range r {
				if fmt.Sprintf("%v", reflect.TypeOf(rr.to)) == childName {
					fss, c, j, o := d.extractField(t.Field(i).Interface())
					m2 := d.models[fmt.Sprintf("%v", reflect.TypeOf(rr.to))]
					join.WriteString(" JOIN " + m2.table + " ON " + fmt.Sprintf("%s.%s", m.table, rr.pK) + "=" + fmt.Sprintf("%s.%s ", m2.table, rr.fK))
					order.WriteString(fmt.Sprintf("%s.%s", m.table, rr.pK) + ",")
					for _, f := range fss {
						fs = append(fs, f)
					}
					if c != "" {
						col.WriteString(c + ",")
						join.WriteString(j)
						order.WriteString(o)
					}
				}
			}
		} else if t.Type().Field(i).Tag.Get("db") != "" {
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

	joinRes := join.String()

	orderRes := order.String()
	if orderRes != "" {
		orderRes = orderRes[:len(orderRes)-1]
	}

	return s, colRes, joinRes, orderRes
}

func (d DBMap) extractField(data interface{}) ([]reflect.StructField, string, string, string) {
	fs := []reflect.StructField{}

	t := reflect.ValueOf(data)

	var col bytes.Buffer
	var join bytes.Buffer
	var order bytes.Buffer

	m := d.models[fmt.Sprintf("%v", reflect.TypeOf(data))]
	r := d.relations[fmt.Sprintf("%v", reflect.TypeOf(data))]

	if len(r) != 0 {
		for _, rr := range r {
			m2 := d.models[fmt.Sprintf("%v", reflect.TypeOf(rr.to))]
			join.WriteString(" JOIN " + m2.table + " ON " + fmt.Sprintf("%s.%s", m.table, rr.pK) + "=" + fmt.Sprintf("%s.%s ", m2.table, rr.fK))
		}
	}

	if m == nil {
		m = d.models[fmt.Sprintf("%v", reflect.TypeOf(data).Elem())]
		r = d.relations[fmt.Sprintf("%v", reflect.TypeOf(data).Elem())]
		t = reflect.New(reflect.TypeOf(data).Elem()).Elem()
		if reflect.TypeOf(t.Interface()).Kind() == reflect.Ptr {
			t = reflect.New(reflect.TypeOf(data).Elem().Elem()).Elem()
			m = d.models[fmt.Sprintf("%v", reflect.TypeOf(data).Elem().Elem())]
			r = d.relations[fmt.Sprintf("%v", reflect.TypeOf(data).Elem().Elem())]
		}
	}

	for i := 0; i < t.NumField(); i++ {
		if t.Type().Field(i).Tag.Get("join") != "" {
			childName := fmt.Sprintf("%v", reflect.TypeOf(t.Field(i).Interface()))
			childName = strings.Replace(childName, "[]", "", 1)
			childName = strings.Replace(childName, "*", "", 1)
			for _, rr := range r {
				if fmt.Sprintf("%v", reflect.TypeOf(rr.to)) == childName {
					m2 := d.models[fmt.Sprintf("%v", reflect.TypeOf(rr.to))]
					join.WriteString(" JOIN " + m2.table + " ON " + fmt.Sprintf("%s.%s", m.table, rr.pK) + "=" + fmt.Sprintf("%s.%s ", m2.table, rr.fK))
					order.WriteString(fmt.Sprintf("%s.%s", m.table, rr.pK) + ",")
					fss, c, j, o := d.extractField(t.Field(i).Interface())
					for _, f := range fss {
						fs = append(fs, f)
					}
					if c != "" {
						col.WriteString(c + ",")
						join.WriteString(j)
						order.WriteString(o)
					}
				}
			}
		} else if t.Type().Field(i).Tag.Get("db") != "" {
			sf := reflect.StructField{
				Name: fmt.Sprintf("%s%s", upperFirst(m.table), t.Type().Field(i).Name),
				Type: t.Type().Field(i).Type,
				Tag:  reflect.StructTag(fmt.Sprintf(`db:"%s.%s"`, m.table, t.Type().Field(i).Tag.Get("db"))),
			}
			col.WriteString(sf.Tag.Get("db") + " as " + `"` + sf.Tag.Get("db") + `"` + ",")
			fs = append(fs, sf)
		}
	}

	colRes := col.String()
	colRes = colRes[:len(colRes)-1]

	joinRes := join.String()

	orderRes := order.String()

	return fs, colRes, joinRes, orderRes
}
