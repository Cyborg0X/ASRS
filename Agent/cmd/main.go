package main

import (
	"bufio"
	"fmt"
	"os"
	"./ASRS/Agent/internal/pkg/checker"
	
)

func main() {
	checkdone := depcheck()
	if checkdone == true {
		buffer := bufio.NewReader(os.Stdin)
		fmt.Printf("Please Enter the IP address of the workstation: ")
		IPaddr, _ := buffer.ReadString('\n')
		fmt.Println(IPaddr)
	}
}
