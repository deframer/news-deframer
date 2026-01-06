package main

import (
	"fmt"
)

func main() {
	fmt.Printf("I am the worker\n")
	select {}
}
