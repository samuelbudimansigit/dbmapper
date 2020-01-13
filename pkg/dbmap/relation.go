package dbmap

import (
	"errors"
	"fmt"
	"reflect"
)

func (d DBMap) AddRelation(table1 interface{}, relationType int, table2 interface{}, tag string, primary interface{}, foreign interface{}) error {
	tableName1 := fmt.Sprintf("%v", reflect.TypeOf(table1))
	tableName2 := fmt.Sprintf("%v", reflect.TypeOf(table2))

	m1 := d.models[tableName1]
	m2 := d.models[tableName2]

	if m1 == nil {
		return errors.New("table not yet registered")
	}

	if m2 == nil {
		return errors.New("table not yet registered")
	}

	var columns1 []Column
	var columns2 []Column

	for _, col := range m1.columns {
		c := Column{
			name: col.name,
		}
		columns1 = append(columns1, c)
	}

	for _, col := range m2.columns {
		c := Column{
			name: col.name,
		}
		columns2 = append(columns2, c)
	}

	m1.data = table1

	m1.columns = columns1
	m2.columns = columns2

	relation := &Relation{
		from:         table1,
		to:           table2,
		pK:           primary,
		fK:           foreign,
		by:           tag,
		relationType: relationType,
	}

	tableName := fmt.Sprintf("%v", reflect.TypeOf(table1))
	d.relations[tableName] = append(d.relations[tableName], relation)

	return nil
}
