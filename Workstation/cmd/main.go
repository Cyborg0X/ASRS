package main

import (
	"fmt"
	"time"

	"github.com/Cyborg0X/ASRS/Workstation/internal/pkg/checker"
	//"github.com/Cyborg0X/ASRS/Workstation/internal/pkg/communication"
	//"github.com/Cyborg0X/ASRS/Workstation/internal/pkg/config"
	"github.com/Cyborg0X/ASRS/Workstation/internal/pkg/handler"
)

func main() {
	checkdone := checker.Depcheck()
	if checkdone == true {
		fmt.Println("Welcome to the Workstation of ASRS")
	}
	//ip, port := config.AgentInfoParser()
	//connection, err := communication.WS_dailer(ip, port)
	//if err != nil {
	//	fmt.Println("connection lost...")
///	}
	handler.TaskHandler()

	time.Sleep(time.Hour)

}
