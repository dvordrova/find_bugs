package main

import (
	"context"
	"fmt"
	"time"

	"github.com/dvordrova/find_bugs/synctest/context_timeout_without_wall_clock/internal/leases"
)

func main() {
	manager := leases.NewManager(5 * time.Second)
	lease := manager.Start(context.Background(), "order-42")
	defer lease.Stop()

	fmt.Printf("lease %s active immediately: %v\n", lease.ID, lease.Active())
	fmt.Println("deadline check is covered by testing/synctest without waiting on wall-clock time")
}
