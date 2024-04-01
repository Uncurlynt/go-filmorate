package models

import (
	"encoding/json"
	"log"
	"time"
)

type CustomTime struct {
	time.Time
}

type ResponseError struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func (t CustomTime) MarshalJSON() ([]byte, error) {
	if t.IsZero() {
		return []byte("null"), nil
	} else {
		return []byte(`"` + t.Format("2006-01-02") + `"`), nil
	}
}

func (t *CustomTime) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		t.Time = time.Time{}
		return nil
	}
	var timeStr string
	if err := json.Unmarshal(data, &timeStr); err != nil {
		return err
	}
	log.Println("Unmarshalled time string:", timeStr)
	parsedTime, err := time.Parse("2006-01-02", timeStr)
	if err != nil {
		return err
	}
	t.Time = parsedTime
	return nil
}
