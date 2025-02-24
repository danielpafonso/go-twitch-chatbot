package twitch

import (
	"testing"
	"time"
)

func TestRateEnd(t *testing.T) {
	// limit to 20 message per 30 seconds
	limiter := RateLimiter{
		End:         time.Now().Add(-time.Second),
		NumberLimit: 20,
		TimeLimit:   30,
		Count:       2,
	}
	if limiter.CanRequest() != true {
		t.Fatal()
	}
	if limiter.Count != 1 {
		t.Fatal("count != than 1")
	}
}

func TestUnderLimit(t *testing.T) {
	// limit to 20 message per 30 seconds
	limiter := RateLimiter{
		End:         time.Now().Add(time.Second),
		NumberLimit: 20,
		TimeLimit:   30,
		Count:       1,
	}
	if limiter.CanRequest() != true {
		t.Fatal()
	}
	if limiter.Count != 2 {
		t.Fatal()
	}
}

func TestLimitReach(t *testing.T) {
	// limit to 20 message per 30 seconds
	limiter := RateLimiter{
		End:         time.Now().Add(time.Second),
		NumberLimit: 20,
		TimeLimit:   30,
		Count:       20,
	}
	// test when limit is reached
	if limiter.CanRequest() != false {
		t.Fatal()
	}
	// test passing the limit
	if limiter.CanRequest() != false {
		t.Fatal()
	}
}

func TestReachingLimit(t *testing.T) {
	// limit to 20 message per 30 seconds
	limiter := RateLimiter{
		End:         time.Now().Add(time.Minute),
		NumberLimit: 20,
		TimeLimit:   30,
	}
	// Under limit
	for i := 0; i < 20; i++ {
		request := limiter.CanRequest()
		if request != true || limiter.Count != i+1 {
			t.Fatal()
		}
	}
	// limit and above
	for i := 0; i < 5; i++ {
		request := limiter.CanRequest()
		if request != false || limiter.Count != limiter.NumberLimit {
			t.Fatal()
		}
	}
}
