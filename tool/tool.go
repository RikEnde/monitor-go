package main

import (
	"fmt"
	u "monitor/util"
	"os"
)

func main() {
	for _, arg := range os.Args[1:] {
		fmt.Printf("%s - %s\n", arg, u.Hash(arg))
	}
}
