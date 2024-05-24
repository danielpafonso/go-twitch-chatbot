package twitch

import (
	"time"
)

type RateLimiter struct {
	Count       int
	End         time.Time
	NumberLimit int
	TimeLimit   int
}

func InitRateLimiter(messageLimit, timeLimit int) RateLimiter {
	return RateLimiter{
		Count: 0,
		End:   time.Now(),
		// End:         time.Now().Add(time.Second * time.Duration(-timeLimit)),
		NumberLimit: messageLimit,
		TimeLimit:   timeLimit,
	}
}

func (rate *RateLimiter) CanRequest() bool {
	// reset limit if necessary
	if rate.End.Before(time.Now()) {
		rate.End = time.Now().Add(time.Second * time.Duration(rate.TimeLimit))
		rate.Count = 1
		return true
	}
	if rate.Count < rate.NumberLimit {
		// under rate limit
		rate.Count += 1
		return true
	} else {
		// rate limit reached
		return false
	}
}
