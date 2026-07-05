package domain

import "time"

type Subscription struct {
	AccountID string
	RenewsAt  time.Time
}

func (s Subscription) NeedsRenewalNotice(window time.Duration) bool {
	remaining := s.RenewsAt.Sub(time.Now())
	return remaining > 0 && remaining <= window
}
