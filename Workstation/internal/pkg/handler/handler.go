package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/Cyborg0X/ASRS/Workstation/internal/pkg/communication"
	"github.com/Cyborg0X/ASRS/Workstation/internal/pkg/config"
)

var red = "\033[31m"
var green = "\033[32m"
var reset = "\033[0m"
var cyan = "\033[36m"

type DataType int

const (
	Typemessage DataType = iota
	TypeSSH
)

type A1A2 struct {
	A string `json:"procedure"`
}

type SSH struct {
	Proceduree string `json:"procedure"`
	Username   string `json:"SSH username"`
}

type DataWrapper struct {
	Type DataType    `json:"type"`
	Data interface{} `json:"data"`
}

func TaskHandler(wgd *sync.WaitGroup) {
	fmt.Println(green + "TASK HANDLER RUNNING NOW" + reset)
	defer wgd.Done()
	var wg sync.WaitGroup
	wg.Add(1)
	get_done := make(chan bool)
	defer close(get_done)
	sendSSH()
	go Get_Status(&wg, get_done)
	go checkIDS()
	wg.Wait()

}

func Get_Status(wg *sync.WaitGroup, getdone <-chan bool) {
	defer wg.Done()
	for {
		time.Sleep(time.Second * 20)
		select {
		case <-getdone:
			return
		default:
			go procedureSelector("A1")
		}
	}
}

func procedureSelector(procedurename string) {
	// pass struct here and get the procedure function name like A1, A2

	switch procedurename {
	case "A1":
		procedure := A1A2{
			A: "A1",
		}
		wrapper := DataWrapper{
			Type: Typemessage,
			Data: procedure,
		}

		jsondata, err := json.MarshalIndent(wrapper, "", "  ")
		if err != nil {
			fmt.Println(red+"\nJSON MESSAGE: Failed to marshal A1:"+reset, err)
			break
		}
		_, err = ProcedureSender(jsondata, procedurename)

		if err != nil {
			fmt.Println(red+"COMMUNICATION MESSAGE: Failed to send A1 to Agent:"+reset, err)
			break
		}

		// handle response logging and shit
	case "A2":

		procedure := A1A2{
			A: "A2",
		}
		wrapper := DataWrapper{
			Type: Typemessage,
			Data: procedure,
		}
		Jsondata, err := json.MarshalIndent(wrapper, "", "  ")
		if err != nil {
			fmt.Println(red+"\nJSON MESSAGE: Failed to marshal A2:"+reset, err)
			return
		}
		_, err = ProcedureSender(Jsondata, procedurename)
		if err != nil {
			fmt.Println(red+"\nCOMMUNICATION MESSAGE: Failed to send A2 to Agent:"+reset, err)
			return
		}

	case "user":
		username := SSHusername()
		fmt.Println(username)
		receiveddata, err := ProcedureSender(username, procedurename)
		SaveKeys(receiveddata)
		if err != nil {
			fmt.Println(red+"\nSSH MESSAGE: Failed to send SSH username to Agent:"+reset, err)
			return
		}
		fmt.Println(green + "\nSSH MESSAGE: SSH username sent successfully" + reset)

	}

}

func SaveKeys(received string) {
	keyspath := "~/.ssh/authorized_keys"
	err := ioutil.WriteFile(keyspath, []byte(received), 0766)
	if err != nil {
		fmt.Println(red+"\nSSH MESSAGE: Failed to write SSH keys:"+reset, err)
	}
}

func ProcedureSender(procedure []byte, procedurename string) (data string, err error) {

	var ip, port = config.AgentInfoParser()

	// Try to send the data, and reconnect if the connection is lost
	conn, err := communication.WS_dailer(ip, port)
	if err != nil {
		conn.Close()
	}
	//defer conn.Close()

	_, err = conn.Write(procedure)
	if err == nil {
		fmt.Printf(green+"\nCOMMUNICATION MESSAGE: Procedure %v sent successfully"+reset, procedurename)
		// response, err := ProcedureReceiver(conn, procedure)
		received := ProcedureReceiver(conn)
		if received == "" {
			fmt.Printf(red+"\nCOMMUNICATION MESSAGE: can't receive after the %v\n%v"+reset, procedurename, err)
			conn.Close()
		}
		return received, nil
	}
	fmt.Println(red+"\nCOMMUNICATION MESSAGE: Failed to send data to agent:"+reset, err)
	return "error", err

}

func sendSSH() {

	filepath := "/etc/ASRS_WS/.config/config.json"
	file, err := ioutil.ReadFile(filepath)
	if err != nil {
		fmt.Println(red+"\nSSH MESSAGE: error reading config file"+reset, err)
	}
	var send config.Config
	err = json.Unmarshal(file, &send)
	fmt.Println(send.Workstationinfo.SendSSH)
	if !send.Workstationinfo.SendSSH {
		procedureSelector("user")
		send.Workstationinfo.SendSSH = true
		jsondata, _ := json.MarshalIndent(send, "", "  ")
		err = ioutil.WriteFile(filepath, jsondata, 0766)
		fmt.Println(green + "\nSSH MESSAGE: SSH has been newly configured in the config file" + reset)

		return
	}
	return
}

func SSHusername() []byte {

	user := exec.Command("whoami")
	output, err := user.Output()
	if err != nil {
		fmt.Println(red + "\nSSH MESSAGE: can't get the username of the device" + reset)
	}
	sshdone := strings.TrimSpace(string(output))
	procedure := SSH{
		Proceduree: "user",
		Username:   sshdone,
	}
	wrapper := DataWrapper{
		Type: TypeSSH,
		Data: procedure,
	}
	jsondata, err := json.MarshalIndent(wrapper, "", "  ")
	if err != nil {
		fmt.Println(red + "\nSSH MESSAGE: can't get the username of the device 2" + reset)
	}
	return jsondata

}

// sent procedure in connection
// receive response in another connection
// in the both ways
// which means big receiver and big sender in each WS and Agent

func ProcedureReceiver(conn net.Conn) (Response string) {

	receive := make([]byte, 1024)
	n, err := conn.Read(receive)
	if err != nil {
		fmt.Println(red+"\nCOMMUNICATION MESSAGE: Failed to read from connection:"+reset, err)
		return ""
	}
	completed := string(receive[:n])
	return completed

	/*
		receivedata := make([]byte, 1024)
		for {
			Numofbytes, err := conn.Read(receivedata)
			if err != nil {
				fmt.Println("\nCOMMUNICATION MESSAGE: Failed to read from connection:", err)
				return "", err
			}
			data := string(receivedata[:Numofbytes])
			if procedure == "A1" {
				return data, nil
			} else if procedure == "A2" {
				// waiting to receive B3
				return procedure, nil
			}
			break
		}
		return
	*/

}

func A1ResponseHandler(response string) {

}
func A2ResponseHandler(response string) {

}

func B3ResquestHandler(reqeust string) {

}

func checkIDS() {
	// if log file detected a attack then procedureSelector
}
