package rulemancer

import (
	"bytes"
	"encoding/json"
	"errors"
)

func DecodeOneOf[T any, U any](c *Config, data []byte, t1 *T, t2 *U) (any, error) {
	targets := []any{t1, t2}
	for _, t := range targets {
		dec := json.NewDecoder(bytes.NewReader(data))
		dec.DisallowUnknownFields()

		if err := dec.Decode(t); err == nil {
			return t, nil
		}
	}
	return nil, errors.New("no matching schema")
}
