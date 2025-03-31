package main

import (
	"testing"
)

func TestInitiate(t *testing.T) {
	filter := contains{}
	pattern := "regex pattern"
	err := filter.Initiate(pattern)
	// assert
	if err != nil {
		t.FailNow()
	}
	if filter.filter != pattern {
		t.FailNow()
	}
	if filter.GetPattern() != pattern {
		t.FailNow()
	}
}

func TestApply(t *testing.T) {
	filter := contains{}
	filter.Initiate("pattern")

	if ok, _ := filter.Apply("Is is a pattern inside a string"); ok != true {
		t.FailNow()
	}
	if ok, _ := filter.Apply("pattern is a word"); ok != true {
		t.FailNow()
	}
	if ok, _ := filter.Apply("This string should fail"); ok != false {
		t.FailNow()
	}
	if ok, _ := filter.Apply("This string with Pattern should fail"); ok != false {
		t.FailNow()
	}
}
