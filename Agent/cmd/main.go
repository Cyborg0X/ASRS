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

	go Program(wtcom,errorchan, eventchan,notificationchan,progresschan)
	<-wtcom

	 
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	// Create the progress bar widget for system progress
	progresswidgets := widgets.NewList()
	progresswidgets.Title = "System Progress"
	progresswidgets.SetRect(0, 0, 50, 10)
	progresswidgets.TextStyle = ui.NewStyle(ui.ColorCyan)

	// Create the list widget for events
	eventWidget := widgets.NewList()
	eventWidget.Title = "Events"
	eventWidget.SetRect(0, 10, 50, 20)
	eventWidget.TextStyle = ui.NewStyle(ui.ColorMagenta)

	// Create the list widget for errors
	errorWidget := widgets.NewList()
	errorWidget.Title = "Errors"
	errorWidget.SetRect(50, 0, 100, 10)
	errorWidget.TextStyle = ui.NewStyle(ui.ColorRed)

	// Create the list widget for notifications
	notificationWidget := widgets.NewList()
	notificationWidget.Title = "Notifications"
	notificationWidget.SetRect(50, 10, 100, 20)
	notificationWidget.TextStyle = ui.NewStyle(ui.ColorGreen)

	ui.Render(progresswidgets, eventWidget, errorWidget, notificationWidget)

	// Simulate system progress and update the progress bar
	/*
		go func() {
			for i := 0; i <= 100; i++ {
				progressBar.Percent = i
				ui.Render(progressBar)
				time.Sleep(100 * time.Millisecond)
			}
		}()
	*/
	// Simulate adding new events, errors, and notifications

	progresscounter := 0
	eventCounter := 0
	errorCounter := 0
	notificationCounter := 0


	go func() {
		for {
			proge := <- progresschan
			progresscounter++
			progresswidgets.Rows = append(progresswidgets.Rows, proge+string(rune(progresscounter)))
			if len(eventWidget.Rows) > 10 {
				progresswidgets.Rows = progresswidgets.Rows[1:]
			}
			ui.Render(progresswidgets)
			time.Sleep(time.Millisecond * 100)
		}
	}()

	go func() {
		for {
			eve := <- eventchan
			eventCounter++
			eventWidget.Rows = append(eventWidget.Rows, eve+string(rune(eventCounter)))
			if len(eventWidget.Rows) > 10 {
				eventWidget.Rows = eventWidget.Rows[1:]
			}
			ui.Render(eventWidget)
			time.Sleep(time.Millisecond * 100)
		}
	}()

	go func() {

		for {
			erro := <- errorchan
			errorCounter++
			errorWidget.Rows = append(errorWidget.Rows, erro+string(rune(errorCounter)))
			if len(errorWidget.Rows) > 10 {
				errorWidget.Rows = errorWidget.Rows[1:]
			}
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
			ui.Render(notificationWidget)
			time.Sleep(time.Millisecond * 100)
		}
	}()
	// Start the main event loop
	uiEvents := ui.PollEvents()
	for {
		select {
		case e := <-uiEvents:
			if e.Type == ui.KeyboardEvent && e.ID == "q" || e.ID == "Q" {
				return
			}
		}
	}
	
}


func Program(kek chan bool, er, eve, noti, prog chan string) {
	var wg sync.WaitGroup
	wg.Add(1)
	checkdone := checker.Depcheck()
	if checkdone {
		fmt.Println(green + "Welcome agent of ASRS" + reset)
	}
	ip, port := handler.WSInfoParser(er)

	procedure_chan := make(chan net.Conn, 1)
	go communication.AG_Listener(ip, port, procedure_chan, eve,er,prog)
	go handler.TaskHandler(&wg, procedure_chan,er,eve,noti,prog)
	wg.Wait()

	// return procedure from listener passed in channel and
	// choose the right procedure and run it
	// let the procedure waiter running too

}


