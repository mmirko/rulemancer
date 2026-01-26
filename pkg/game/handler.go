package game

import (
	"errors"
	"io"
	"net/http"
)

func GenericHandler(c *Config, w http.ResponseWriter, r *http.Request) (string, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return "", err
	}

	var move Move

	switch v, err := DecodeOneOf(c, body, &move); {
	case err != nil:
		if c.Debug {
			println("Failed to decode request payload:", string(body))
		}
		return "", errors.New("invalid request payload")
	case v == &move:
		return move2ClipsAssert(move), nil
	default:
		return "", errors.New("unknown request payload")
	}
}
