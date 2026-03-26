package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

type BaseModel struct {
	CreatedAt time.Time // 创建时间
	UpdatedAt time.Time // 更新时间
}

type JSONMap map[string]any

func (m JSONMap) Value() (driver.Value, error) {
	if m == nil {
		return nil, nil
	}
	return json.Marshal(m)
}

func (m *JSONMap) Scan(value any) error {
	if value == nil {
		*m = make(map[string]any)
		return nil
	}
	var err error
	switch v := value.(type) {
	case []byte:
		err = json.Unmarshal(v, m)
	case string:
		err = json.Unmarshal([]byte(v), m)
	default:
		err = errors.New("basetypes.JSONMap.Scan: invalid value type")
	}
	if err != nil {
		return err
	}
	return nil
}
