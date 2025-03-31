package main

import (
	"errors"
	"strings"
)

type contains struct {
	filter string
}

func (flt *contains) Initiate(args ...interface{}) error {
	if len(args) < 1 {
		return errors.New("no string to match")
	}

	flt.filter = args[0].(string)
	return nil
}

func (flt *contains) GetPattern() string {
	return flt.filter
}

func (flt *contains) Apply(line string) (bool, error) {
	ok := strings.Contains(line, flt.filter)
	return ok, nil
}

// Export Symbol
var Filter contains
