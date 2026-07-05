package main

import (
	"fmt"

	"github.com/dvordrova/find_bugs/concurrency/select_priority_assumption/internal/dispatcher"
)

func main() {
	fmt.Println("select does not prioritize the first ready case")
	fmt.Println("run make lint to see the repeated schedule check")

	high := make(chan dispatcher.Job, 1)
	low := make(chan dispatcher.Job, 1)
	high <- dispatcher.Job{Queue: "high", ID: "urgent-1"}
	low <- dispatcher.Job{Queue: "low", ID: "batch-1"}

	job := dispatcher.Next(high, low)
	fmt.Printf("one run selected %s job %s\n", job.Queue, job.ID)
}
