package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/Cyborg0X/ASRS/Workstation/internal/pkg/checker"
	//"github.com/Cyborg0X/ASRS/Workstation/internal/pkg/communication"
	//"github.com/Cyborg0X/ASRS/Workstation/internal/pkg/config"
	"github.com/Cyborg0X/ASRS/Workstation/internal/pkg/handler"

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

				ui.Render(progressWidget, eventWidget, errorWidget, notificationWidget)
			} else if e.Type == ui.KeyboardEvent && (e.ID == "<Down>" || e.ID == "j") {

				eventWidget.ScrollDown()

			}
		}
	}
}








func Program(kek chan bool, er, eve, noti, prog chan string) {
	
	var wg sync.WaitGroup
	wg.Add(1)
	checkdone := checker.Depcheck()
	if checkdone {
		fmt.Println(green+"Welcome to the Workstation of ASRS"+reset)
		kek <- true
	}
	//ip, port := config.AgentInfoParser()
	//connection, err := communication.WS_dailer(ip, port)
	//if err != nil {
	//	fmt.Println("connection lost...")
	///	}
	
	go handler.TaskHandler(&wg, noti,er, eve, prog)
	wg.Wait()

}
