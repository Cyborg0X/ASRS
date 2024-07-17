package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"sync"
	"time"

	"github.com/Cyborg0X/ASRS/Workstation/internal/pkg/communication"
	"github.com/Cyborg0X/ASRS/Workstation/internal/pkg/config"
)

var red = "\033[31m"
var green = "\033[32m"
var reset = "\033[0m"

type DataType int

const (
	TypeA1 DataType = iota
	TypeA2
)

type A1 struct {
	A string `json:"procedure"`
}

type A2 struct {
	A          string `json:"procedure"`
	AttackerIP string `json:"Attacker IP"`
	AttackTime string `json:"Time of attack"`
}

/*
	type SSH struct {
		Proceduree string `json:"procedure"`
		Username   string `json:"SSH username"`
		Pass       string `json:"SSH pass"`
	}
*/
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
	go Get_Status(&wg, get_done)

	go checkIDS()

	wg.Wait()

}

func checkIDS() {
	// if log file detected a attack then procedureSelector
	// check IDS 
	//.......

	// SEND procedure if detected an attack
	slp := make(chan bool, 1)
	for {
		procedureSelector("A2", slp)
		if !<-slp {
			break
		}
		time.Sleep(time.Minute * 1)
	}

}

func Get_Status(wg *sync.WaitGroup, getdone <-chan bool) {
	defer wg.Done()
	for {
		time.Sleep(time.Second * 20)
		select {
		case <-getdone:
			return
		default:
			sleep := make(chan bool, 1)
			go procedureSelector("A1", sleep); if <-sleep {time.Sleep(time.Minute * 3)}
		}
	}
}

func procedureSelector(procedurename string, slp chan bool) {
	// pass struct here and get the procedure function name like A1, A2
	
	switch procedurename {
	case "A1":
		procedure := A1{
			A: "A1",
		}
		wrapper := DataWrapper{
			Type: TypeA1,
			Data: procedure,
		}

		jsondata, err := json.MarshalIndent(wrapper, "", "  ")
		if err != nil {
			fmt.Println(red+"\nJSON MESSAGE: Failed to marshal A1:"+reset, err)
			break
		}
		receivedmessage, err := ProcedureSender(jsondata, procedurename)
		if err != nil {
			fmt.Println(red+"\nCOMMUNICATION MESSAGE: Failed to send A1 to Agent:"+reset, err)
			return
		}
		if string(receivedmessage) == "B3PROC" {
			slp <- true
			return
		} else {
			slp <- false
			return
		}

		// handle response logging and shit
	case "A2":

		var config config.Config
		filepath := "/etc/ASRS_WS/.config/config.json"
		file, err := ioutil.ReadFile(filepath)
		if err != nil {
			fmt.Println(red+"\nSelector MESSAGE: Error reading config file: %w"+reset, err) // Wrap error with context for other errors
		}
		err = json.Unmarshal(file, &config)
		if err != nil {
			fmt.Println(red+"\nSelector MESSAGE: Error Unmarshal config file: %w"+reset, err) // Wrap error with context for other errors
		}
		IP := config.Detectionmarker.AttackerIP
		time := config.Detectionmarker.AttackTiming

		procedure := A2{
			A:	"A2",
			AttackerIP: IP,
			AttackTime: time,

		}
		wrapper := DataWrapper{
			Type: TypeA2,
			Data: procedure,
		}
		Jsondata, err := json.MarshalIndent(wrapper, "", "  ")
		if err != nil {
			fmt.Println(red+"\nJSON MESSAGE: Failed to Marshal A2:"+reset, err)
			return
		}
		receivedmessage, err := ProcedureSender(Jsondata, procedurename)
		if err != nil {
			fmt.Println(red+"\nCOMMUNICATION MESSAGE: Failed to send A2 to Agent:"+reset, err)
			return
		}
		if string(receivedmessage) == "B3PROC" {
			slp <- true
			return
		} else {
			slp <- false
			return
		}

		/*
			case "user":
				username := SSHusername()
				fmt.Println(string(username))
				receiveddata, err := ProcedureSender(username, procedurename)
				SaveKeys(receiveddata)
				if err != nil {
					fmt.Println(red+"\nSSH MESSAGE: Failed to send SSH username to Agent:"+reset, err)
					return
				}
				fmt.Println(green + "\nSSH MESSAGE: SSH username sent successfully" + reset)


		*/
	}
}

/*
func SaveKeys(received []byte) {
	keyspath := "/root/.ssh/authorized_keys"
	err := ioutil.WriteFile(keyspath, received, 0766)
	if err != nil {
		fmt.Println(red+"\nSSH MESSAGE: Failed to write SSH keys:"+reset, err)
	}
}
*/

func ProcedureSender(procedure []byte, procedurename string) (data []byte, err error) {

	ip, port := config.AgentInfoParser()

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
		received, err := ProcedureReceiver(conn)
		if err != nil {
			fmt.Printf(red+"\nCOMMUNICATION MESSAGE: can't receive after the %v\n%v"+reset, procedurename, err)

		}
		return received, nil
	}
	fmt.Println(red+"\nCOMMUNICATION MESSAGE: Failed to send data to agent:"+reset, err)
	return
}

/*
func sendSSH() {

	filepath := "/etc/ASRS_WS/.config/config.json"
	file, err := ioutil.ReadFile(filepath)
	if err != nil {
		fmt.Println(red+"\nSSH MESSAGE: error reading config file"+reset, err)
	}
	var send config.Config
	err = json.Unmarshal(file, &send)
	if err != nil {
		fmt.Println("SSH MESSAGE: Failed to Unmarshal for SSH USERNAME")
	}
	fmt.Println(send.Workstationinfo.SendSSH)
	if !send.Workstationinfo.SendSSH {
		procedureSelector("user")
		send.Workstationinfo.SendSSH = true
		jsondata, _ := json.MarshalIndent(send, "", "  ")
		err = ioutil.WriteFile(filepath, jsondata, 0766)
		fmt.Println(green + "\nSSH MESSAGE: SSH has been newly configured in the config file" + reset)
	}
}


func SSHusername() []byte {
	filepath := "/etc/ASRS_WS/.config/config.json"
	file, err := ioutil.ReadFile(filepath)
	if err != nil {
		fmt.Println(red+"\nSSH MESSAGE: error reading config file"+reset, err)
	}
	var send config.Config
	err = json.Unmarshal(file, &send)
	if err != nil {
		fmt.Println("SSH MESSAGE: Failed to Unmarshal for SSH PASSWORD")
	}

	//////
	user := exec.Command("whoami")
	output, err := user.Output()
	if err != nil {
		fmt.Println(red + "\nSSH MESSAGE: can't get the username of the device" + reset)
	}

	username := strings.TrimSpace(string(output))
	pass := send.Workstationinfo.SSHpass
	procedure := SSH{
		Proceduree: "user",
		Username:   username,
		Pass:       pass,
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
*/

// sent procedure in connection
// receive response in another connection
// in the both ways
// which means big receiver and big sender in each WS and Agent

func ProcedureReceiver(conn net.Conn) (Response []byte, err error) {

	buffer := make([]byte, 1024)
	receveddata := bytes.Buffer{}
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Println(red+"\nCOMMUNICATION MESSAGE: Failed to read from connection:"+reset, err)
		}
		receveddata.Write(buffer[:n])
		if n < len(buffer) {
			break
		}
	}
	data := receveddata.Bytes()
	return data, err

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
