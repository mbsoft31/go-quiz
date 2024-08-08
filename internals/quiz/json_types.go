package quiz

import (
	"database/sql/driver"
	"encoding/json"
	"gorm.io/gorm"
)

// JSONMap is a type to handle map[string]interface{} fields
type JSONMap map[string]interface{}

// Scan implements the Scanner interface for JSONMap
func (jm *JSONMap) Scan(value interface{}) error {
	if value == nil {
		*jm = JSONMap{}
		return nil
	}
	data, ok := value.([]byte)
	if !ok {
		return gorm.ErrInvalidData
	}
	return json.Unmarshal(data, jm)
}

// Value implements the Valuer interface for JSONMap
func (jm JSONMap) Value() (driver.Value, error) {
	if len(jm) == 0 {
		return nil, nil
	}
	return json.Marshal(jm)
}
