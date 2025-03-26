package main

import (
	"fmt"
	"time"
)

type uptime struct {
	start time.Time
}

func (cmd *uptime) Initiate(args ...interface{}) error {
	cmd.start = time.Now()

	return nil
}

func (cmd *uptime) Execute(args ...interface{}) (string, error) {
	duration := time.Since(cmd.start)

	return fmt.Sprint(duration), nil
}

// Export Symbol
var Command uptime
