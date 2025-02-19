package service

import (
	"encoding/json"
	"fmt"
)

type Typeable struct {
	Type string `json:"type"`
}

func TypeOf(data []byte) (string, error) {
	var t Typeable

	err := json.Unmarshal(data, &t)
	if err != nil {
		return "", fmt.Errorf("unmarshaling: %w", err)
	}

	return t.Type, nil
}
