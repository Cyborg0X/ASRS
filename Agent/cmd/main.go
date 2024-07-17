package main

import (
	"fmt"
	"net"
	"sync"

	"github.com/Cyborg0X/ASRS/Agent/internal/pkg/checker"
	"github.com/Cyborg0X/ASRS/Agent/internal/pkg/communication"
	"github.com/Cyborg0X/ASRS/Agent/internal/pkg/handler"
)
var green = "\033[32m"
var reset = "\033[0m"


func main() {
	var wg sync.WaitGroup
	wg.Add(1)
	checkdone := checker.Depcheck()
	if checkdone == true {
		fmt.Println(green+"Welcome agent of ASRS"+reset)
	}
	ip, port := handler.WSInfoParser()
	
	procedure_chan := make(chan net.Conn, 1)
	go communication.AG_Listener(ip, port, procedure_chan)
	go handler.TaskHandler(&wg, procedure_chan)
	wg.Wait()

	// return procedure from listener passed in channel and
	// choose the right procedure and run it
	// let the procedure waiter running too

}


