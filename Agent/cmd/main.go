package main

import (
	"fmt"
	"net"
	"sync"

	"github.com/Cyborg0X/ASRS/Agent/internal/pkg/checker"
	"github.com/Cyborg0X/ASRS/Agent/internal/pkg/communication"
	"github.com/Cyborg0X/ASRS/Agent/internal/pkg/handler"
	"github.com/Cyborg0X/ASRS/Agent/internal/pkg/logger"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(1)
	checkdone := checker.Depcheck()
	if checkdone == true {
		fmt.Println("Welcome agent of ASRS")
	}
	ip, port := handler.WSInfoParser()
	procedure_chan := make(chan net.Conn, 1)
	go communication.AG_Listener(ip, port, procedure_chan)
	B3 := logger.DetectionMarker()
	go handler.TaskHandler(&wg, procedure_chan, B3)
	wg.Wait()

	// return procedure from listener passed in channel and
	// choose the right procedure and run it
	// let the procedure waiter running too

}
