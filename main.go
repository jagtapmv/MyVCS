package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: myvcs <command>")
		return
	}

	command := os.Args[1]

	fmt.Printf("Your entered command is: %s\n", command)
}
