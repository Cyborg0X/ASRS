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

func AG_Listener(ip string, port string, channel chan net.Conn) error {
	//retryDelay := 5 * time.Second
	var listenerr net.Listener
	var err error
	var connf net.Conn
	address := fmt.Sprintf(":%v", port)
	for {

		listenerr, err = net.Listen("tcp", address)
		if err != nil {
			fmt.Println("COMMUNICATION MESSAGE: Failed to create a listener", err)
			fmt.Println(address)
			continue
		}

		for {
			connf, err = listenerr.Accept()
			if err != nil {
				fmt.Println("COMMUNICATION MESSAGE: Failed to accept connection", err)
				continue
			} else {
				fmt.Println("COMMUNICATION MESSAGE: Connection Accepted")
				channel <- connf
			}

		}
	}
}

func AG_dialer() {

}

func Response_Sender(message string, conn net.Conn) {
	for {
		_, err := conn.Write([]byte(message))
		if err != nil {
			fmt.Println("Faild to send ", message)
			continue
		}
	}
}

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

func errorhandler(err error, s string) {
	if err != nil {
		fmt.Println("Error: ", s, err)
	}
}

func AssignWorkstationIP() error {
	filepath := "/etc/ASRS_agent/.config/config.json"

	file, err := ioutil.ReadFile(filepath)
	if err != nil {

		return fmt.Errorf("Error reading config file: %w", err) // Wrap error with context for other errors
	}
	errorhandler(err, "Error reading config file: ")

	var lookforip handler.Config
	err = json.Unmarshal(file, &lookforip)
	errorhandler(err, "Error parsing config file: ")

	// if IP address of the agent and workstation is found then continue without entering the IPs

	fmt.Printf("Do you want to change or set Workstation IP address? N/y : ")
	changeip := bufio.NewReader(os.Stdin)
	input, _ := changeip.ReadString('\n')
	g := strings.TrimSpace(input)
	switch g {
	case "Y", "y":
		fmt.Printf("Enter the IP address of the Workstation: ")
		buffer := bufio.NewReader(os.Stdin)
		lookforip.Workstationinfo.IPaddr, _ = buffer.ReadString('\n')
		lookforip.Workstationinfo.IPaddr = strings.TrimSpace(lookforip.Workstationinfo.IPaddr) // Remove trailing newline
		modifiedData, err := json.MarshalIndent(lookforip, "", "  ")
		errorhandler(err, "Error marshaling JSON: ")
		err = ioutil.WriteFile(filepath, modifiedData, 0766)
		errorhandler(err, "Error writing JSON file:")
		fmt.Printf("Your Agent IP address is: %v\n", lookforip.Workstationinfo.IPaddr)
		fmt.Printf("Your Agent port is: %v\n", lookforip.Workstationinfo.Port)
	case "n", "N":
		fmt.Printf("Your Agent IP address is: %v\n", lookforip.Workstationinfo.IPaddr)
		fmt.Printf("Your Agent port is: %v\n", lookforip.Workstationinfo.Port)
	}

	return nil
}

func AssignAgentIP() error {
	filepath := "/etc/ASRS_agent/.config/config.json"

	file, err := ioutil.ReadFile(filepath)
	if err != nil {

		return fmt.Errorf("Error reading config file: %w", err) // Wrap error with context for other errors
	}
	errorhandler(err, "Error reading config file: ")

	var lookforip handler.Config
	err = json.Unmarshal(file, &lookforip)
	errorhandler(err, "Error parsing config file: ")

	// if IP address of the agent and workstation is found then continue without entering the IPs

	fmt.Printf("Do you want to change or set Agent IP address? N/y : ")
	changeip := bufio.NewReader(os.Stdin)
	input, _ := changeip.ReadString('\n')
	g := strings.TrimSpace(input)
	switch g {
	case "Y", "y":
		fmt.Printf("Enter the IP address of the Agent: ")
		buffer := bufio.NewReader(os.Stdin)
		lookforip.Agentinfo.Ipaddr, _ = buffer.ReadString('\n')
		lookforip.Agentinfo.Ipaddr = strings.TrimSpace(lookforip.Agentinfo.Ipaddr) // Remove trailing newline
		modifiedData, err := json.MarshalIndent(lookforip, "", "  ")
		errorhandler(err, "Error marshaling JSON: ")
		err = ioutil.WriteFile(filepath, modifiedData, 0766)
		errorhandler(err, "Error writing JSON file:")
		fmt.Printf("Your Agent IP address is: %v\n", lookforip.Agentinfo.Ipaddr)
		fmt.Printf("Your Agent port is: %v\n", lookforip.Agentinfo.Port)
	case "n", "N":
		fmt.Printf("Your Agent IP address is: %v\n", lookforip.Agentinfo.Ipaddr)
		fmt.Printf("Your Agent port is: %v\n", lookforip.Agentinfo.Port)
	}

	return nil
}

//func connectToWS()  {
//	connection, err := net.D}
