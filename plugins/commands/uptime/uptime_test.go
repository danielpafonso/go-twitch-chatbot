package main

import (
	"strings"
	"testing"
	"time"
)

func TestInitiate(t *testing.T) {
	cmd := uptime{}

	before := time.Now()
	err := cmd.Initiate()

	// assert
	if err != nil {
		t.Fatal("error not nil")
	}
	if before.Before(cmd.start) != true {
		t.Fatal("uptime not greather than before")
	}
	if time.Now().After(cmd.start) != true {
		t.Fatal("uptime not lower than after")
	}
}

func TestExecute(t *testing.T) {
	cmd := uptime{
		start: time.Now().Add(-1 * time.Second),
	}
	output, err := cmd.Execute()

	// assert
	if err != nil {
		t.Fatal("error not nil")
	}
	if strings.HasPrefix(output, "1.0") != true {
		t.Fatalf("expected '1.0' prefix, got: %s", output)
	}
	if strings.HasSuffix(output, "s") != true {
		t.Fatalf("expected 's' sufix, got: %s", output)
	}
}
