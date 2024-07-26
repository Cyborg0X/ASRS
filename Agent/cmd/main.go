
package main

import (

	"log"
	"time"
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)


var progresschan = make(chan string)
var eventchan = make(chan string)
var errorchan = make(chan string)
var notificationchan = make(chan string)

var progress = []string{"Agent Listner STARTED", 
						"ASRS TASK HANDLER STARTED",
						"LOCAL ACTIONS STARTED","ASRS RESTORE BACKUP STARTED", 
						"COMMUNICATION MESSAGE: Connection from workstation Accepted",
						}

var noti = []string{"A2 PROCEDURE RECEIVED",
					"CHECKING FOR ATTACKER IP ADDRESS",
					"ATTACKER IP ADDRESS FOUND",
					"STARTING CHECKING FOR BREACHED FILES",
					"CURRENT SYSTEM AND THE BACKUP ARE NOT THE SAME",
					"STARTING SELF-HEALING PROCESS",
					"RESTORING THE SYSTEM TO PRE-ATTACK STATE",
					"ALL CHANGES HAVE BEEN REVENTED",
					}




var events = []string{"RSYNC MESSAGE: Sending status to workstation","ASRS initiating syncing process...", 
					"RSYNC MESSAGE: syncing home files and web files to workstation\nit might take a long time.......", 
					"RSYNC MESSAGE: SYSTEM AND WEB FILES SYNCED",
				}
var errors = []string{""}



func main() {

	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()
	progresscounter := 0
	eventCounter := 0
	errorCounter := 0
	notificationCounter := 0
	// Create the widgets
	progressWidget := widgets.NewList()
	progressWidget.Title = "Agent System Progress"
	progressWidget.SetRect(0, 0, 50, 10)
	progressWidget.TextStyle = ui.NewStyle(ui.ColorCyan)
	progressWidget.WrapText = true

	eventWidget := widgets.NewList()
	eventWidget.Title = "Agent Events"
	eventWidget.SetRect(0, 10, 50, 20)
	eventWidget.TextStyle = ui.NewStyle(ui.ColorMagenta)
	eventWidget.WrapText = true

	errorWidget := widgets.NewList()
	errorWidget.Title = "Agent Errors"
	errorWidget.SetRect(50, 0, 100, 10)
	errorWidget.TextStyle = ui.NewStyle(ui.ColorRed)
	errorWidget.WrapText = true

	notificationWidget := widgets.NewList()
	notificationWidget.Title = "Agent Notifications"
	notificationWidget.SetRect(50, 10, 100, 20)
	notificationWidget.TextStyle = ui.NewStyle(ui.ColorGreen)
	notificationWidget.WrapText = true

	ui.Render(progressWidget, eventWidget, errorWidget, notificationWidget)


	go func() {
		for {
			proge := <-progresschan
			progresscounter++
			progressWidget.Rows = append(progressWidget.Rows, proge+string(rune(progresscounter)))
			if len(progressWidget.Rows) > 10 {
				progressWidget.Rows = progressWidget.Rows[1:]
			}
			time.Sleep(time.Second *1)
			ui.Render(progressWidget)
			time.Sleep(time.Millisecond * 100)
		}
	}()

	go func() {
		for {
			eve := <-eventchan
			eventCounter++
			eventWidget.Rows = append(eventWidget.Rows, eve+string(rune(eventCounter)))
			if len(eventWidget.Rows) > 10 {
				eventWidget.Rows = eventWidget.Rows[1:]
			}
			time.Sleep(time.Second *1)
			ui.Render(eventWidget)
			time.Sleep(time.Millisecond * 100)
		}
	}()

	go func() {

		for {
			erro := <-errorchan
			errorCounter++
			errorWidget.Rows = append(errorWidget.Rows, erro+string(rune(errorCounter)))
			if len(errorWidget.Rows) > 10 {
				errorWidget.Rows = errorWidget.Rows[1:]
			}
			time.Sleep(time.Second *1)
			ui.Render(errorWidget)
			time.Sleep(time.Millisecond * 100)
		}
	}()

	go func() {
		for {
			noti := <-notificationchan
			notificationCounter++
			notificationWidget.Rows = append(notificationWidget.Rows, noti+string(rune(notificationCounter)))
			if len(notificationWidget.Rows) > 10 {
				notificationWidget.Rows = notificationWidget.Rows[1:]
			}
			time.Sleep(time.Second *1)
			ui.Render(notificationWidget)
			time.Sleep(time.Millisecond * 100)
		}
	}()
	for i := 0; i < len(progress); i++ {
		progresschan<- progress[i]
	}

	for i := 0; i < len(events); i++ {
		eventchan<- events[i]
	}

	for i := 0; i < len(noti); i++ {
		notificationchan<- noti[i]
	}

	for i := 0; i < len(errors); i++ {
		errorchan<- errors[i]
	}


	uiEvents := ui.PollEvents()
	for {
		select {
		case e := <-uiEvents:
			if e.Type == ui.KeyboardEvent && (e.ID == "q" || e.ID == "Q" || e.ID == "<C-c>") {
				return
			}
			// Handle scrolling events
			if e.Type == ui.KeyboardEvent && (e.ID == "<Up>" || e.ID == "k") {
				eventWidget.ScrollUp()
				ui.Render(eventWidget)
			} else if e.Type == ui.KeyboardEvent && (e.ID == "<Down>" || e.ID == "j") {
				eventWidget.ScrollDown()
				ui.Render(eventWidget)
			}
		}
	}

	
} 







/*
[**] [1:1000023:0] Command Injection detected - /vulnerabilities/exec [**]
[Classification: Attempted User Privilege Gain] [Priority: 1]
07/26/2024-10:30:15.123456 TCP 192.168.1.100:56789 -> 203.0.113.1:80
*/















/*
package main

import (
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"github.com/Cyborg0X/ASRS/Agent/internal/pkg/checker"
	"github.com/Cyborg0X/ASRS/Agent/internal/pkg/communication"
	"github.com/Cyborg0X/ASRS/Agent/internal/pkg/handler"
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

var green = "\033[32m"
var reset = "\033[0m"
var progresschan = make(chan string)
var eventchan = make(chan string)
var errorchan = make(chan string)
var notificationchan = make(chan string)

func main() {

	wtcom := make(chan bool)
	go Program(wtcom, errorchan, eventchan, notificationchan, progresschan)
	<-wtcom

	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()
	progresscounter := 0
	eventCounter := 0
	errorCounter := 0
	notificationCounter := 0
	// Create the widgets
	progressWidget := widgets.NewList()
	progressWidget.Title = "System Progress"
	progressWidget.SetRect(0, 0, 50, 10)
	progressWidget.TextStyle = ui.NewStyle(ui.ColorCyan)
	progressWidget.WrapText = true

	eventWidget := widgets.NewList()
	eventWidget.Title = "Events"
	eventWidget.SetRect(0, 10, 50, 20)
	eventWidget.TextStyle = ui.NewStyle(ui.ColorMagenta)
	eventWidget.WrapText = true

	errorWidget := widgets.NewList()
	errorWidget.Title = "Errors"
	errorWidget.SetRect(50, 0, 100, 10)
	errorWidget.TextStyle = ui.NewStyle(ui.ColorRed)
	errorWidget.WrapText = true

	notificationWidget := widgets.NewList()
	notificationWidget.Title = "Notifications"
	notificationWidget.SetRect(50, 10, 100, 20)
	notificationWidget.TextStyle = ui.NewStyle(ui.ColorGreen)
	notificationWidget.WrapText = true

	ui.Render(progressWidget, eventWidget, errorWidget, notificationWidget)


	go func() {
		for {
			proge := <-progresschan
			progresscounter++
			progressWidget.Rows = append(progressWidget.Rows, proge+string(rune(progresscounter)))
			if len(progressWidget.Rows) > 10 {
				progressWidget.Rows = progressWidget.Rows[1:]
			}
			time.Sleep(time.Second *1)
			ui.Render(progressWidget)
			time.Sleep(time.Millisecond * 100)
		}
	}()

	go func() {
		for {
			eve := <-eventchan
			eventCounter++
			eventWidget.Rows = append(eventWidget.Rows, eve+string(rune(eventCounter)))
			if len(eventWidget.Rows) > 10 {
				eventWidget.Rows = eventWidget.Rows[1:]
			}
			time.Sleep(time.Second *1)
			ui.Render(eventWidget)
			time.Sleep(time.Millisecond * 100)
		}
	}()

	go func() {

		for {
			erro := <-errorchan
			errorCounter++
			errorWidget.Rows = append(errorWidget.Rows, erro+string(rune(errorCounter)))
			if len(errorWidget.Rows) > 10 {
				errorWidget.Rows = errorWidget.Rows[1:]
			}
			time.Sleep(time.Second *1)
			ui.Render(errorWidget)
			time.Sleep(time.Millisecond * 100)
		}
	}()

	go func() {
		for {
			noti := <-notificationchan
			notificationCounter++
			notificationWidget.Rows = append(notificationWidget.Rows, noti+string(rune(notificationCounter)))
			if len(notificationWidget.Rows) > 10 {
				notificationWidget.Rows = notificationWidget.Rows[1:]
			}
			time.Sleep(time.Second *1)
			ui.Render(notificationWidget)
			time.Sleep(time.Millisecond * 100)
		}
	}()


	uiEvents := ui.PollEvents()
	for {
		select {
		case e := <-uiEvents:
			if e.Type == ui.KeyboardEvent && (e.ID == "q" || e.ID == "Q" || e.ID == "<C-c>") {
				return
			}
			// Handle scrolling events
			if e.Type == ui.KeyboardEvent && (e.ID == "<Up>" || e.ID == "k") {
				eventWidget.ScrollUp()
				ui.Render(eventWidget)
			} else if e.Type == ui.KeyboardEvent && (e.ID == "<Down>" || e.ID == "j") {
				eventWidget.ScrollDown()
				ui.Render(eventWidget)
			}
		}
	}
} 

//
func Program(kek chan bool, er, eve, noti, prog chan string) {
	var wg sync.WaitGroup
	wg.Add(1)
	checkdone := checker.Depcheck()
	if checkdone {
		fmt.Println(green + "Welcome agent of ASRS" + reset)
		kek <- true
	}
	ip, port := handler.WSInfoParser(er)

	procedure_chan := make(chan net.Conn, 1)
	go communication.AG_Listener(ip, port, procedure_chan, eve, er, prog)
	go handler.TaskHandler(&wg, procedure_chan, er, eve, noti, prog)
	wg.Wait()

	// return procedure from listener passed in channel and
	// choose the right procedure and run it
	// let the procedure waiter running too

}
*/