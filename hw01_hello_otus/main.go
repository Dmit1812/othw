package main

import (
	"fmt"

	"golang.org/x/example/stringutil"
)

func main() {
	hellowOtus := "Hello, OTUS!"
	fmt.Printf("%s", stringutil.Reverse(hellowOtus))
}
