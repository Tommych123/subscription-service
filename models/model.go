package models

import (
	"database/sql/driver"
	"fmt"
	"time"
)

// swagger:model MonthYear
type MonthYear struct {
	time.Time
}

func (MonthYear) SwaggerType() []string {
	return []string{"string"}
}

func (MonthYear) SwaggerFormat() string {
	return "month-year"
}

const monthYearFormat = "01-2006"

func (m *MonthYear) UnmarshalJSON(data []byte) error {
	str := string(data)
	if len(str) >= 2 && str[0] == '"' && str[len(str)-1] == '"' {
		str = str[1 : len(str)-1]
	}
	t, err := time.Parse(monthYearFormat, str)
	if err != nil {
		return fmt.Errorf("invalid month-year format: %v", err)
	}
	m.Time = t
	return nil
}

func (m MonthYear) MarshalJSON() ([]byte, error) {
	return []byte(`"` + m.Format(monthYearFormat) + `"`), nil
}

func (m *MonthYear) Scan(value interface{}) error {
	if value == nil {
		m.Time = time.Time{}
		return nil
	}
	switch v := value.(type) {
	case time.Time:
		m.Time = v
		return nil
	case []byte:
		t, err := time.Parse("2006-01-02", string(v))
		if err != nil {
			return err
		}
		m.Time = t
		return nil
	case string:
		t, err := time.Parse("2006-01-02", v)
		if err != nil {
			return err
		}
		m.Time = t
		return nil
	default:
		return fmt.Errorf("cannot scan type %T into MonthYear", value)
	}
}

func (m MonthYear) Value() (driver.Value, error) {
	if m.Time.IsZero() {
		return nil, nil
	}
	return m.Time.Format("2006-01-02"), nil
}

type Subscription struct {
	ID          string     `db:"id" json:"id"`
	ServiceName string     `db:"service_name" json:"service_name" example:"Spotify"`
	Price       int        `db:"price" json:"price" example:"9"`
	UserID      string     `db:"user_id" json:"user_id" example:"550e8400-e29b-41d4-a716-446655440000"`
	StartDate   MonthYear  `db:"start_date" json:"start_date" swaggertype:"string"`
	EndDate     *MonthYear `db:"end_date" json:"end_date" swaggertype:"string"`
}
