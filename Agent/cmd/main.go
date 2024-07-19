package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"sync"

	"github.com/Cyborg0X/ASRS/Agent/internal/pkg/checker"
	"github.com/Cyborg0X/ASRS/Agent/internal/pkg/communication"
	"github.com/Cyborg0X/ASRS/Agent/internal/pkg/handler"
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

var green = "\033[32m"
var reset = "\033[0m"

func main() {
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

	//progresscounter := 0
	//eventCounter := 0
	//errorCounter := 0
	//notificationCounter := 0
	//progresschan := make(chan string)
	//eventchan := make(chan string)
	//errorchan := make(chan string)
	//notificationchan := make(chan string)

	Program()

	go func() {
		// progress events from file
		progContent, err := readFile("progress.txt")
		if err != nil {
			log.Printf("failed to read events file: %v", err)
		} else {
			eventWidget.Rows = progContent
		}
		ui.Render(progresswidgets)

	}()

	go func() {
		// Read events from file
		eventContent, err := readFile("events.txt")
		if err != nil {
			log.Printf("failed to read events file: %v", err)
		} else {
			eventWidget.Rows = eventContent
		}

		ui.Render(eventWidget)
	}()

	go func() {
			// Read errors from file
			errorContent, err := readFile("errors.txt")
			if err != nil {
				log.Printf("failed to read errors file: %v", err)
			} else {
				errorWidget.Rows = errorContent
			}
		ui.Render(errorWidget)
	}()

	go func() {
			// Read notifications from file
			notificationContent, err := readFile("notifications.txt")
			if err != nil {
				log.Printf("failed to read notifications file: %v", err)
			} else {
				notificationWidget.Rows = notificationContent
			}
		ui.Render(notificationWidget)
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

func Program() {
	var wg sync.WaitGroup
	wg.Add(1)
	checkdone := checker.Depcheck()
	if checkdone {
		fmt.Println(green + "Welcome agent of ASRS" + reset)
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

func readFile(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var content []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		content = append(content, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return content, nil
}
