package main

import (
	"fmt"
	"time"

	"github.com/dvordrova/find_bugs/teamrules/no_wall_clock_in_domain/internal/billing/domain"
)

func main() {
	subscription := domain.Subscription{
		AccountID: "tenant-42",
		RenewsAt:  time.Now().Add(24 * time.Hour),
	}

	fmt.Printf("renewal notice due: %v\n", subscription.NeedsRenewalNotice(48*time.Hour))
}
