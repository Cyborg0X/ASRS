package communication

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"strings"

	"github.com/Cyborg0X/ASRS/Agent/internal/pkg/handler"
)
var red = "\033[31m"
var green = "\033[32m"
var reset = "\033[0m"
var cyan = "\033[36m"

func AG_Listener(ip string, port string, channel chan net.Conn, eve, erro, prog chan string) error {
	handler.EventHandler("AG LISTNER STARTED", eve)
	//retryDelay := 5 * time.Second
	var listenerr net.Listener
	var err error
	var connf net.Conn
	address := fmt.Sprintf(":%v", port)
	for {

		listenerr, err = net.Listen("tcp", address)
		if err != nil {
			handler.Errorhandler(err,"COMMUNICATION: Failed to create a listener", erro)
			continue
		}

		for {
			connf, err = listenerr.Accept()
			if err != nil {
				handler.Errorhandler(err,"COMMUNICATION MESSAGE: Failed to accept connection", erro)
				continue
			} else {
				handler.ProgHandler("COMMUNICATION MESSAGE: Connection Accepted", prog)
				channel <- connf
			}

		}
	}
}

func AG_dialer() {

}
/*
func Response_Sender(message string, conn net.Conn) {
	fmt.Println("RESPONSE SENDER STARTED")
	for {
		_, err := conn.Write([]byte(message))
		if err != nil {
			fmt.Println(red+"SENDER MESSAGE: Faild to send"+reset, message)
			continue
		}
	}
}
*/
func Connect_to_ws(ipaddr string, port string) {
	wsinfo := handler.Config{}
	wsinfo.Workstationinfo.IPaddr = ipaddr
	wsinfo.Workstationinfo.Port = port
	//json_wsinfo, err := json.Marshal(wsinfo)
	//if err != nil {
	//	fmt.Println("Error:", err)
	//	return
	//}
	//infofile, err := os.Open("/etc/ASRS_agent/.config/config.json")
}



func AssignWorkstationIP() error {
	filepath := "/etc/ASRS_agent/.config/config.json"

	file, err := ioutil.ReadFile(filepath)
	if err != nil {
		fmt.Println("error reading config file",err)
		//return //fmt.Errorf(red+"error reading config file: %v"+reset, err) // Wrap error with context for other errors
	}
	//errorhandler(err, "Error reading config file: ")

	var lookforip handler.Config
	err = json.Unmarshal(file, &lookforip)
	if err != nil {
	fmt.Println(red+"Error parsing config file: "+reset,err)
	}
	// if IP address of the agent and workstation is found then continue without entering the IPs

	fmt.Printf(cyan+"Do you want to change or set Workstation IP address? N/y : "+reset)
	changeip := bufio.NewReader(os.Stdin)
	input, _ := changeip.ReadString('\n')
	g := strings.TrimSpace(input)
	switch g {
	case "Y", "y":
		fmt.Printf(cyan+"Enter the IP address of the Workstation: "+reset)
		buffer := bufio.NewReader(os.Stdin)
		lookforip.Workstationinfo.IPaddr, _ = buffer.ReadString('\n')
		lookforip.Workstationinfo.IPaddr = strings.TrimSpace(lookforip.Workstationinfo.IPaddr) // Remove trailing newline
		modifiedData, err := json.MarshalIndent(lookforip, "", "  ")
		fmt.Println(err, red+"Error marshaling JSON: "+reset)
		fileper, _ := os.Stat(filepath)
		per := fileper.Mode().Perm()
		err = ioutil.WriteFile(filepath, modifiedData, per)
		fmt.Println(err, red+"Error writing JSON file:"+reset)
		fmt.Printf(cyan+"Your Agent IP address is: %v\n"+reset, lookforip.Workstationinfo.IPaddr)
		fmt.Printf(cyan+"Your Agent port is: %v\n"+reset, lookforip.Workstationinfo.Port)
	case "n", "N":
		fmt.Printf(cyan+"Your Agent IP address is: %v\n"+reset, lookforip.Workstationinfo.IPaddr)
		fmt.Printf(cyan+"Your Agent port is: %v\n"+reset, lookforip.Workstationinfo.Port)
	}

	return nil
}

func AssignAgentIP() error {
	filepath := "/etc/ASRS_agent/.config/config.json"

	file, err := ioutil.ReadFile(filepath)
	if err != nil {

		return fmt.Errorf(red+"Error reading config file: %w"+reset, err) // Wrap error with context for other errors
	}
	fmt.Println(err, red+"Error reading config file: "+reset)

	var lookforip handler.Config
	err = json.Unmarshal(file, &lookforip)
	fmt.Println(err, red+"Error parsing config file: "+reset)

	// if IP address of the agent and workstation is found then continue without entering the IPs

	fmt.Printf(cyan+"Do you want to change or set Agent IP address? N/y : "+reset)
	changeip := bufio.NewReader(os.Stdin)
	input, _ := changeip.ReadString('\n')
	g := strings.TrimSpace(input)
	switch g {
	case "Y", "y":
		fmt.Printf(cyan+"Enter the IP address of the Agent: "+reset)
		buffer := bufio.NewReader(os.Stdin)
		lookforip.Agentinfo.Ipaddr, _ = buffer.ReadString('\n')
		lookforip.Agentinfo.Ipaddr = strings.TrimSpace(lookforip.Agentinfo.Ipaddr) // Remove trailing newline
		modifiedData, err := json.MarshalIndent(lookforip, "", "  ")
		fmt.Println(err, red+"Error marshaling JSON: "+reset)
		fileper, _ := os.Stat(filepath)
		per := fileper.Mode().Perm()
		err = ioutil.WriteFile(filepath, modifiedData, per)
		fmt.Println(err, red+"Error writing JSON file:"+reset)
		fmt.Printf(cyan+"Your Agent IP address is: %v\n"+reset, lookforip.Agentinfo.Ipaddr)
		fmt.Printf(cyan+"Your Agent port is: %v\n"+reset, lookforip.Agentinfo.Port)
	case "n", "N":
		fmt.Printf(cyan+"Your Agent IP address is: %v\n"+reset, lookforip.Agentinfo.Ipaddr)
		fmt.Printf(cyan+"Your Agent port is: %v\n"+reset, lookforip.Agentinfo.Port)
	}

	return nil
}

//func connectToWS()  {
//	connection, err := net.D}
