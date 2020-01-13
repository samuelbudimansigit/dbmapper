package dbmap

import (
	"fmt"
	"reflect"
)

func (d DBMap) MapRecursive(data interface{}, values []map[string]interface{}, flag int) {
	dt := reflect.ValueOf(data)
	res := reflect.ValueOf(data).Elem()
	var dd interface{}
	var m *Model
	l := 0

	if reflect.TypeOf(data).Kind() == reflect.Struct {
		dt = reflect.ValueOf(data)
		m = d.models[fmt.Sprintf("%v", reflect.TypeOf(data))]
	}
	if reflect.TypeOf(data).Kind() == reflect.Slice {
		dd = reflect.New(reflect.TypeOf(data).Elem().Elem()).Interface()
		dt = reflect.ValueOf(dd).Elem()
		m = d.models[fmt.Sprintf("%v", reflect.TypeOf(dd).Elem())]
	}
	if reflect.TypeOf(data).Kind() == reflect.Ptr {
		dd = reflect.New(reflect.TypeOf(data).Elem().Elem()).Interface()
		dt = reflect.ValueOf(dd).Elem()
		m = d.models[fmt.Sprintf("%v", reflect.TypeOf(data).Elem().Elem())]
	}
	for i := 0; i < flag; i++ {
		val := values[i]
		for j := 0; j < dt.NumField(); j++ {
			v := val[fmt.Sprintf("%s.%s", m.table, dt.Type().Field(j).Tag.Get("db"))]
			if v != nil {
				dt.Field(j).Set(reflect.ValueOf(v))
			} else if dt.Type().Field(j).Tag.Get("join") != "" {
				l = d.MapRecursiveChild(data, values, j, i, dt, *m, flag)
			}
		}
		i = l
		res.Set(reflect.Append(res, dt))
		dt = reflect.New(reflect.TypeOf(dt.Interface())).Elem()
	}
}

func (d DBMap) MapRecursiveChild(data interface{}, values []map[string]interface{}, flag int, curr int, field reflect.Value, parent Model, total int) int {
	var m *Model
	isPtr := 0
	dt := reflect.ValueOf(data).Elem()
	dat := reflect.ValueOf(data).Elem()
	typeSlice := 1
	var dd interface{}

	dat = field.Field(flag)

	dt = dat
	dd = dt.Interface()
	last := curr + 1
	l := 0

	if reflect.TypeOf(dd).Kind() == reflect.Struct {
		dt = reflect.ValueOf(dd)
		m = d.models[fmt.Sprintf("%v", reflect.TypeOf(dd))]
	}
	if reflect.TypeOf(dd).Kind() == reflect.Slice {
		dt = reflect.New(reflect.TypeOf(reflect.ValueOf(dd).Interface()).Elem()).Elem()
		m = d.models[fmt.Sprintf("%v", reflect.TypeOf(dd).Elem())]
		typeSlice = 2
		dat = reflect.New(reflect.TypeOf(dat.Interface()).Elem()).Elem()
		if reflect.TypeOf(dat.Interface()).Kind() == reflect.Ptr {
			m = d.models[fmt.Sprintf("%v", reflect.TypeOf(dd).Elem().Elem())]
			dat = reflect.New(reflect.TypeOf(dat.Interface()).Elem()).Elem()
			dt = reflect.New(reflect.TypeOf(reflect.ValueOf(dd).Interface()).Elem().Elem()).Elem()
			isPtr = 1
		}
	}

	if reflect.TypeOf(dd).Kind() == reflect.Ptr {
		dt = reflect.New(reflect.TypeOf(dd).Elem()).Elem()
		dat = reflect.New(reflect.TypeOf(dat.Interface()).Elem()).Elem()
		m = d.models[fmt.Sprintf("%v", reflect.TypeOf(dd).Elem())]
		isPtr = 1
	}

	for i := curr; i < total; i++ {
		val := values[i]
		for j := 0; j < dt.NumField(); j++ {
			v := val[fmt.Sprintf("%s.%s", m.table, dt.Type().Field(j).Tag.Get("db"))]
			// log.Println(v)
			// log.Println(fmt.Sprintf("%s.%s", m.table, dt.Type().Field(j).Tag.Get("db")))
			if v != nil {
				dat.Field(j).Set(reflect.ValueOf(v))
			} else if dt.Type().Field(j).Tag.Get("join") != "" {
				l = d.MapRecursiveChild(data, values, j, i, dat, *m, total)
				i = l
			}
		}
		after := values[i+1]
		f := 0
		for i := 0; i < len(parent.primaryKey); i++ {
			if val[fmt.Sprintf("%s.%s", parent.table, parent.primaryKey[i])] == after[fmt.Sprintf("%s.%s", parent.table, parent.primaryKey[i])] {
				f++
			}
		}
		if f != len(parent.primaryKey) {
			last = i
			if typeSlice == 2 {
				if isPtr == 1 {
					field.Field(flag).Set(reflect.Append(field.Field(flag), dat.Addr()))
				} else {
					field.Field(flag).Set(reflect.Append(field.Field(flag), dat))
				}
			} else {
				if isPtr == 1 {
					field.Field(flag).Set(dat.Addr())
				} else {
					field.Field(flag).Set(dat)
				}
			}
			break
		} else {
			if typeSlice == 2 {
				if isPtr == 1 {
					field.Field(flag).Set(reflect.Append(field.Field(flag), dat.Addr()))
				} else {
					field.Field(flag).Set(reflect.Append(field.Field(flag), dat))
				}
			} else {
				if isPtr == 1 {
					field.Field(flag).Set(dat.Addr())
				} else {
					field.Field(flag).Set(dat)
				}
			}
		}
		dat = reflect.New(reflect.TypeOf(dat.Interface())).Elem()
	}
	return last
}
