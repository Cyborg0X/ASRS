package main

import (
	"fmt"
	"github.com/Cyborg0X/ASRS/Agent/internal/pkg/checker"
)

func main() {
	checkdone := checker.Depcheck()
	if checkdone == true {
		fmt.Println("Welcome")
	}
}
