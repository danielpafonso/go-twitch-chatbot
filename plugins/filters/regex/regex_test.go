package main

import (
	"testing"
)

func TestInitiate(t *testing.T) {
	filter := regex{}
	err := filter.Initiate(".+")
	if err != nil {
		t.FailNow()
	}
	if filter.pattern != filter.GetPattern() {
		t.FailNow()
	}
}

func TestApply(t *testing.T) {
	var tests = []struct {
		name     string
		pattern  string
		messages []string
		expected []bool
	}{
		{
			"or match",
			"[a|b]d",
			[]string{
				"ads are bad",
				"commics are called bds",
				"this should fail",
			},
			[]bool{
				true,
				true,
				false,
			},
		},
		{
			"digits keywords",
			"\\d+",
			[]string{
				"I have 4 apples",
				"The year was 1991",
				"this should fail",
			},
			[]bool{
				true,
				true,
				false,
			},
		},
		{
			"only timestamps",
			`^[0-9]{4}-\d{2}-\d\d$`,
			[]string{
				"1991-01-10",
				"this should fail",
				"1991-1991",
				"19911991",
				"This date, 1991-19-91, must fail",
			},
			[]bool{
				true,
				false,
				false,
				false,
				false,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			filter := regex{}
			err := filter.Initiate(test.pattern)
			if err != nil {
				t.FailNow()
			}
			if len(test.messages) != len(test.expected) {
				t.FailNow()
			}
			for i := range test.messages {
				match, _ := filter.Apply(test.messages[i])
				if match != test.expected[i] {
					t.Fail()
				}
			}
		})
	}
}
