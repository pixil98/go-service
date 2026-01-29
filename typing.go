package service

import (
	"encoding/json"
	"errors"
	"fmt"
)

// ErrMissingType is returned when TypeOf is called with data that has no type field set.
var ErrMissingType = errors.New("missing or empty type field")

type Typeable struct {
	Type string `json:"type"`
}

func TypeOf(data []byte) (string, error) {
	var t Typeable

	err := json.Unmarshal(data, &t)
	if err != nil {
		return "", fmt.Errorf("unmarshaling: %w", err)
	}

	if t.Type == "" {
		return "", ErrMissingType
	}

	return t.Type, nil
}
