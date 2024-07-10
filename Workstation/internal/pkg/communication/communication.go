package communication

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"strings"
	"time"

	"github.com/Cyborg0X/ASRS/Workstation/internal/pkg/config"
)

// change this
// listner on the agent
// sender in the workstation
// close connection each time you finish procedure

// if the workstion sends heal_now then also run the B3 listner

// move it to the agent
/*
func WS_Listener(port string, ip string) (net.Conn,error) {
	retryDelay := 5 * time.Second
	var listenerr net.Listener
	var err error
	for{

		listenerr, err = net.Listen("tcp",port)
		if err == nil {
			fmt.Println("Listener created successfully")
			break
		}
		fmt.Printf("Failed to create listner %v\n: Retrying in %v....\n", err, retryDelay)
		time.Sleep(retryDelay)
	}
	for{
		conn, err := listenerr.Accept()
		if err != nil {
			fmt.Println("Failed to accept connection", err)
			continue
		} else {
			fmt.Println("connection accepted successfully")
			return conn, err
		}
	//	fmt.Printf("Failed to accept the connection %v\n: Retrying in %v....\n", err, retryDelay)
	//	time.Sleep(retryDelay)
	}

}
*/
var red = "\033[31m"
var green = "\033[32m"
var reset = "\033[0m"
var cyan = "\033[36m"

func WS_dailer(ip, port string) (net.Conn, error) {
	fmt.Println("WS DAILER STARTED")
	retryDelay := 3 * time.Second
	var err error
	var conn net.Conn
	addre := fmt.Sprintf("%s:%s", ip, port)
	for {
		conn, err = net.Dial("tcp", addre)
		if err == nil {

			fmt.Println(green+"\nCOMMUNICATION MESSAGE: The Workstation connected to the agent successfully...."+reset)
			return conn, err
		}
		fmt.Printf(red+"\nCOMMUNICATION ERROR: Failed to connect to the agent %v\n: Retrying in %v....\n"+reset, err, retryDelay)
		time.Sleep(retryDelay)

	}

}

func Connect_to_ws(ipaddr string, port string) {
	wsinfo := config.Config{}
	wsinfo.Workstationinfo.IPaddr = ipaddr
	wsinfo.Workstationinfo.Port = port
	//json_wsinfo, err := json.Marshal(wsinfo)
	//if err != nil {
	//	fmt.Println("Error:", err)
	//	return
	//}
	//infofile, err := os.Open("/etc/ASRS_WS/.config/config.json")
}

func errorhandler(err error, s string) {
	if err != nil {
		fmt.Println(red+"Error: "+reset, s, err)
	}
}

func AssignWorkstationIP() error {
	filepath := "/etc/ASRS_WS/.config/config.json"

	file, err := ioutil.ReadFile(filepath)
	if err != nil {

		return fmt.Errorf(red+"\nCONFIG ERROR: Error reading config file: %w"+reset, err) // Wrap error with context for other errors
	}
	errorhandler(err, red+"Error reading config file: "+reset)

	var lookforip config.Config
	err = json.Unmarshal(file, &lookforip)
	errorhandler(err, "\nCONFIG ERROR: Error parsing config file")

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
		errorhandler(err, "Error marshaling JSON: ")
		err = ioutil.WriteFile(filepath, modifiedData, 0766)
		errorhandler(err, "Error writing JSON file:")
		fmt.Printf(cyan+"Your Workstation IP address is: %v\n"+reset, lookforip.Workstationinfo.IPaddr)
		fmt.Printf(cyan+"Your Workstation port is: %v\n"+reset, lookforip.Workstationinfo.Port)
	case "n", "N":
		fmt.Printf(cyan+"Your Agent IP address is: %v\n"+reset, lookforip.Workstationinfo.IPaddr)
		fmt.Printf(cyan+"Your Agent port is: %v\n"+reset, lookforip.Workstationinfo.Port)

	}

	return nil
}

func AssignAgentIP() error {
	filepath := "/etc/ASRS_WS/.config/config.json"

	file, err := ioutil.ReadFile(filepath)
	if err != nil {

		return fmt.Errorf(red+"error reading config file: %w"+reset, err) // Wrap error with context for other errors
	}
	errorhandler(err, "Error reading config file: ")

	var lookforip config.Config
	err = json.Unmarshal(file, &lookforip)
	errorhandler(err, "Error parsing config file: ")

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
		errorhandler(err, "Error marshaling JSON: ")
		/////pass 
		fmt.Printf(cyan+"Enter the SSH pass to give to the Agent: "+reset)
		buffer = bufio.NewReader(os.Stdin)
		lookforip.Workstationinfo.SSHpass, _ = buffer.ReadString('\n')
		lookforip.Workstationinfo.SSHpass = strings.TrimSpace(lookforip.Agentinfo.Ipaddr) // Remove trailing newline
		modifiedData, err := json.MarshalIndent(lookforip, "", "  ")
		errorhandler(err, "Error marshaling JSON: ")
		/////
		err = ioutil.WriteFile(filepath, modifiedData, 0766)
		errorhandler(err, "Error writing JSON file:")
		fmt.Printf(cyan+"Your Agent IP address is: %v\n"+reset, lookforip.Agentinfo.Ipaddr)
		fmt.Printf(cyan+"Your Agent port is: %v\n"+reset, lookforip.Agentinfo.Port)
	case "n", "N":
		fmt.Printf(cyan+"Your Agent IP address is: %v\n"+reset, lookforip.Agentinfo.Ipaddr)
		fmt.Printf(cyan+"Your Agent port is: %v\n"+reset, lookforip.Agentinfo.Port)
	}

	return nil
}
