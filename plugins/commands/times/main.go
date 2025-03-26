package main

import (
	"time"
)

type times struct{}

func (cmd *times) Initiate(args ...interface{}) error {
	return nil
}

func (cmd *times) Execute(args ...interface{}) (string, error) {
	stime := time.Now().Format(time.RFC850)
	return stime, nil
}

// Export Symbol
var Command times
