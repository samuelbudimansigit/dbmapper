package types

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type Metadata map[string]interface{}

// Value override value's function for metadata (ADT) type
func (p Metadata) Value() (driver.Value, error) {
	j, err := json.Marshal(p)
	return j, err
}

// Scan override scan's function for metadata (ADT) type
func (p *Metadata) Scan(src interface{}) error {
	source, ok := src.([]byte)
	if !ok {
		return errors.New("Type assertion .([]byte) failed")
	}

	var i interface{}
	err := json.Unmarshal(source, &i)
	if err != nil {
		return err
	}

	*p, ok = i.(map[string]interface{})
	if !ok {
		return errors.New("Type assertion .(map[string]interface{}) failed")
	}

	return nil
}
