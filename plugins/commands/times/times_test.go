package main

import (
	"testing"
	"time"
)

func TestInitiate(t *testing.T) {
	cmd := times{}
	err := cmd.Initiate()
	// assert
	if err != nil {
		t.FailNow()
	}
}

func TestExecute(t *testing.T) {
	cmd := times{}

	before := time.Now().Add(-1 * time.Second)
	stime, err1 := cmd.Execute()
	parsed, err2 := time.Parse(time.RFC850, stime)
	after := time.Now().Add(time.Second)
	// assert
	if err1 != nil {
		t.FailNow()
	}
	if err2 != nil {
		t.FailNow()
	}
	if before.Before(parsed) != true {
		t.Fatalf("before not before\n%s\n%s", before, parsed)
	}
	if after.After(parsed) != true {
		t.Fatal("after not after")
	}
}
