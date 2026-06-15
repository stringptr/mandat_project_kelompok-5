package main

import (
	"fmt"
	"os"

	"si-gizi/backend/internal/hash"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run genhash.go <password>")
		return
	}
	h, err := hash.Hash(os.Args[1])
	if err != nil {
		panic(err)
	}
	fmt.Println(h)
}
