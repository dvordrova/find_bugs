package main

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

const maxCustomerIDLine = 32

func ImportCustomerIDs(r io.Reader) ([]string, error) {
	scanner := bufio.NewScanner(r)
	scanner.Buffer(make([]byte, 0, maxCustomerIDLine), maxCustomerIDLine)

	ids := make([]string, 0, 16)
	for scanner.Scan() {
		ids = append(ids, scanner.Text())
	}

	return ids, nil
}

func main() {
	input := strings.Repeat("customer-", 5) + "42\n"

	ids, err := ImportCustomerIDs(strings.NewReader(input))
	if err != nil {
		fmt.Printf("import failed: %v\n", err)
		return
	}

	fmt.Printf("imported %d customer ids\n", len(ids))
}
