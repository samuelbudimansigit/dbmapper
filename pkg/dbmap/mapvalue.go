package dbmap

import (
	"reflect"
)

func (d DBMap) ValueToMapString(data interface{}) ([]map[string]interface{}, int) {
	res := make([]map[string]interface{}, MAX_RETURN)
	x := reflect.ValueOf(data)

	flag := 0

	for j := 0; j < x.Len(); j++ {
		temp := make(map[string]interface{})
		y := x.Index(j).Elem()
		for i := 0; i < y.NumField(); i++ {
			temp[y.Type().Field(i).Tag.Get("db")] = y.Field(i).Interface()
		}
		res[flag] = temp
		flag++
	}
	return res, flag
}
