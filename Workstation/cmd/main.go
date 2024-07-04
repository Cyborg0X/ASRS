package main

import (
	"fmt"
	"sync"

	"github.com/Cyborg0X/ASRS/Workstation/internal/pkg/checker"
	//"github.com/Cyborg0X/ASRS/Workstation/internal/pkg/communication"
	//"github.com/Cyborg0X/ASRS/Workstation/internal/pkg/config"
	"github.com/Cyborg0X/ASRS/Workstation/internal/pkg/handler"
)

func main() {
	var wg sync.WaitGroup

	checkdone := checker.Depcheck()
	if checkdone == true {
		fmt.Println("Welcome to the Workstation of ASRS")
	}
	//ip, port := config.AgentInfoParser()
	//connection, err := communication.WS_dailer(ip, port)
	//if err != nil {
	//	fmt.Println("connection lost...")
	///	}
	wg.Add(1)
	handler.TaskHandler(&wg)
	wg.Wait()

}
