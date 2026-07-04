package main

import (
	"fmt"

	"github.com/acme/contractsdk"
)

func main() {
	plan := contractsdk.DefaultPlan

	printPlan(plan)
}

func printPlan(plan *contractsdk.Plan) {
	fmt.Printf("tenant uses %s plan\n", plan.Name)
}
