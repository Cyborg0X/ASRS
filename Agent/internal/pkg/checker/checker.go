package checker

import (
	//"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	//"io/ioutil"
	"os/exec"
	"strings"

	//"github.com/Cyborg0X/ASRS/Agent/internal/pkg/communication"
	"github.com/Cyborg0X/ASRS/Agent/internal/pkg/communication"
	"github.com/Cyborg0X/ASRS/Agent/internal/pkg/handler"
)

var red = "\033[31m"
var green = "\033[32m"
var reset = "\033[0m"

func Depcheck() bool {
	// return ack of
	packages := make([]string, 6)
	packages[0] = "rsync"
	packages[1] = "snapper"
	packages[2] = "ssh"
	packages[3] = "snort"
	packages[4] = "openssh-server"
	packages[5] = "openssh-client"
	checklist := make(map[string]string)
	update := exec.Command("sudo", "apt", "update")
	update.Run()
	for i := range packages {
		checkpkg := exec.Command("sudo", "dpkg", "-s", packages[i])
		//install := exec.Command("sudo", "apt", "install", packages[i], "-y")
		output, err := checkpkg.CombinedOutput()
		outputstr := string(output)
		if err != nil {
			checklist[packages[i]] = "NOT installed"
			//fmt.Printf("installing %v..... please wait\n", packages[i])
			//install.Run()

		} else if strings.Contains(outputstr, "Status: install ok installed") {

			checklist[packages[i]] = "installed"

		}
		// don't forget to add some lines to interact with a text-based user interface (TUI)
	}
	var feelsgood int
	for key, value := range checklist {
		//fmt.Printf("%v %v\n", key, value)
		if value == "installed" {
			fmt.Printf(green+"%v is installed\n"+reset, key)
			feelsgood++
		} else {
			fmt.Println(red+"Some packages not intalled!!!!"+reset)
			fmt.Println(key, value)
		} 
	}

	if feelsgood == 6 {
		fmt.Println(green+"ALL PACKAGES HAS BEEN INSTALLED SECCUSSFULLY"+reset)
		fmt.Println(green+"Initilizing config file,,,,,,"+reset)
		filepath := "/etc/ASRS_agent/.config/config.json"

		_, err := ioutil.ReadFile(filepath)
		file, _ := os.Stat(filepath)
		// If file doesn't exist, assume it's the first run and create a new one with InitializeJSON
		if os.IsNotExist(err) || file.Size() == 0 {
			err := handler.InitializeJSON()
			if err != nil {
				fmt.Println(red+"Error initialize config file"+reset)
			}

		}

		err = communication.AssignWorkstationIP()
		if err != nil {
			fmt.Println(red+"Erroring assigning Workstation IP address"+reset)
		}
		err = communication.AssignAgentIP()
		if err != nil {
			fmt.Println(red+"Erroring assigning Agent IP address"+reset)
		}
		var detector handler.Config
		fmt.Println(green+"marker : "+reset, detector.Detectionmarker)

	} else {
		panic(red+"Error checking startup"+reset)
	}
	fmt.Fprint(os.Stdout, "\x1b[H\x1b[2J")
	ack := true
	return ack
}
