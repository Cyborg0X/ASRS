package main

import (
	"fmt"
	"bufio"
	"os"
	"github.com/Cyborg0X/ASRS/Workstation/internal/pkg/checker"

)


func main() {
	checkdone := checker.Depcheck()
	if checkdone == true {
		buffer := bufio.NewReader(os.Stdin)
		fmt.Printf("Please Enter the IP address of the workstation: ")
		IPaddr, _ := buffer.ReadString('\n')
		fmt.Println(IPaddr)
	}
}

