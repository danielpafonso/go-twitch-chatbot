package main

import (
	"errors"
	"regexp"
)

type regex struct {
	pattern  string
	compiled *regexp.Regexp
}

func (flt *regex) Initiate(args ...interface{}) error {
	if len(args) < 1 {
		return errors.New("no pattern string to initiate regex")
	}

	flt.pattern = args[0].(string)
	var err error
	flt.compiled, err = regexp.Compile(flt.pattern)
	if err != nil {
		return err
	}
	return nil
}

func (flt *regex) GetPattern() string {
	return flt.pattern
}

func (flt *regex) Apply(line string) (bool, error) {
	match := flt.compiled.MatchString(line)
	return match, nil
}

// Export Symbol
var Filter regex
