package game

import (
	"errors"
	"io"
	"log"
	"net/http"
	"os"
)

func GenericAssertHandler(c *Config, w http.ResponseWriter, r *http.Request) (string, string, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return "", "", err
	}

	var move Move

	switch v, err := DecodeOneOf(c, body, &move); {
	case err != nil:
		if c.Debug {
			l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, red("[game/GenericAssertHandler]")+" ", 0)
			l.Println("Failed to decode request payload:", string(body))
		}
		return "", "", errors.New("invalid request payload")
	case v == &move:
		return "move", move2ClipsAssert(move), nil
	default:
		if c.Debug {
			l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, red("[game/GenericAssertHandler]")+" ", 0)
			l.Println("Unknown request payload:", string(body))
		}
		return "", "", errors.New("unknown request payload")
	}
}
