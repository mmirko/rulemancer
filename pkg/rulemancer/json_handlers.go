package rulemancer

import (
	"errors"
	"log"
	"os"
)

func jsonGenericDecoder(c *Config, body []byte) ([]string, error) {
	var type1 []map[string][]string

	switch v, err := DecodeOneOf(c, body, &type1); {
	case err != nil:
		if c.Debug {
			l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, red("[rulemancer/jsonGenericDecoder]")+" ", 0)
			l.Println("Failed to decode request payload:", string(body))
		}
		return nil, errors.New("invalid request payload")
	case v == &type1:
		return assertType1(type1), nil
	default:
		if c.Debug {
			l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, red("[rulemancer/jsonGenericDecoder]")+" ", 0)
			l.Println("Unknown request payload:", string(body))
		}
		return nil, errors.New("unknown request payload")
	}
}

func assertType1(in []map[string][]string) []string {
	result := make([]string, 0)
	for _, item := range in {
		fact := ""
		for name, values := range item {
			fact += "(" + name + " "
			for _, value := range values {
				fact += value + " "
			}
			fact += ") "
		}
		result = append(result, fact)
	}
	return result
}
