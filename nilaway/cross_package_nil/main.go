package main

import (
	"fmt"
	"log"

	"github.com/dvordrova/find_bugs/nilaway/cross_package_nil/internal/profile"
)

func main() {
	repo := profile.NewRepository()

	p, err := repo.FindByEmail("nobody@example.com")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("sending welcome email to %s\n", p.Email)
}
