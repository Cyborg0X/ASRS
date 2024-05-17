package communication

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/Cyborg0X/ASRS/Agent/internal/pkg/handler"
)


func Connect_to_ws(ipaddr string, port string)  {
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


func Checkandassign_wsip() error {
	filepath := "/etc/ASRS_agent/.config/config.json"
	
	file, err := ioutil.ReadFile(filepath)
	if err != nil {
        // If file doesn't exist, assume it's the first run and create a new one with InitializeJSON
        if os.IsNotExist(err) {
            err = handler.InitializeJSON()
            if err != nil {
                return err
            }
            return Checkandassign_wsip() // Retry Checkandassign_wsip after creating the file
        }
        return fmt.Errorf("Error reading config file: %w", err) // Wrap error with context for other errors
    }
	errorhandler(err,"Error reading config file: ")
	
	var lookforip handler.Config
	err = json.Unmarshal(file, &lookforip)
	errorhandler(err,"Error parsing config file: ")

// if IP address of the agent and workstation is found then continue without entering the IPs
 
	if lookforip.Workstationinfo.IPaddr == ""  {
		fmt.Printf("Enter the IP address of the Workstation: ")
		buffer := bufio.NewReader(os.Stdin)
		lookforip.Workstationinfo.IPaddr, _ = buffer.ReadString('\n')
		lookforip.Workstationinfo.IPaddr = strings.TrimSpace(lookforip.Workstationinfo.IPaddr) // Remove trailing newline
		modifiedData, err := json.Marshal(lookforip)
    	errorhandler(err,"Error marshaling JSON: ")

		err = ioutil.WriteFile(filepath, modifiedData, 0644)
    	errorhandler(err,"Error writing JSON file:")

	} else {
		fmt.Printf("Your Agent IP address is: %v", lookforip.Workstationinfo.IPaddr)
		fmt.Printf("Your Agent port is: %v", lookforip.Workstationinfo.Port)
	}
	return nil
}