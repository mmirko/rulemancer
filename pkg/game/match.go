package game

import (
	"bytes"
	"encoding/json"
	"errors"
)

type Config struct {
	Debug bool
}

func DecodeOneOf[T any](c *Config, data []byte, targets ...*T) (*T, error) {
	for _, t := range targets {
		dec := json.NewDecoder(bytes.NewReader(data))
		dec.DisallowUnknownFields()

		if err := dec.Decode(t); err == nil {
			return t, nil
		}
	}
	return nil, errors.New("no matching schema")
}
