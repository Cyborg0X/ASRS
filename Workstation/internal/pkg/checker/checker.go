package checker

import (
	//"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	//"io/ioutil"
	"os/exec"
	"strings"

	"github.com/Cyborg0X/ASRS/Workstation/internal/pkg/communication"
	"github.com/Cyborg0X/ASRS/Workstation/internal/pkg/config"
)

var red = "\033[31m"
var green = "\033[32m"
var reset = "\033[0m"
var cyan = "\033[36m"

func Depcheck() bool {
	// return ack of
	packages := make([]string, 2)
	packages[0] = "rsync"
	packages[1] = "openssh-server"

	checklist := make(map[string]string)
	update := exec.Command("sudo", "dnf", "update")
	update.Run()
	for i := range packages {
		checkpkg := exec.Command("rpm", "-qi", packages[i])
		install := exec.Command("sudo", "dnf", "install", packages[i], "-y")
		output, err := checkpkg.CombinedOutput()
		outputstr := string(output)
		if err != nil {
			checklist[packages[i]] = "NOT installed"
			//fmt.Printf("installing %v..... please wait\n", packages[i])
			install.Run()

		} else if strings.Contains(outputstr, "Name") {

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
			fmt.Println(green + "Some packages not intalled!!!!" + reset)
			fmt.Println(key, value)
		}
	}

	if feelsgood == 2 {
		fmt.Println(green + "ALL PACKAGES HAS BEEN INSTALLED SECCUSSFULLY" + reset)
		fmt.Println(green + "Initilizing config file,,,,,," + reset)
		filepath := "/etc/ASRS_WS/.config/config.json"

		_, err := ioutil.ReadFile(filepath)
		file, _ := os.Stat(filepath)
		// If file doesn't exist, assume it's the first run and create a new one with InitializeJSON
		if os.IsNotExist(err) || file.Size() == 0 {
			err := config.InitializeJSON()
			if err != nil {
				fmt.Println(red + "Error initialize config file" + reset)
			}

		}

		err = communication.AssignWorkstationIP()
		if err != nil {
			fmt.Println(red+"Erroring assigning Workstation IP address"+reset, err)
		}
		err = communication.AssignAgentIP()
		if err != nil {
			fmt.Println(red+"Erroring assigning Agent IP address"+reset, err)
		}

	} // else {
	//panic("Error checking startup")
	//	}
	fmt.Fprint(os.Stdout, "\x1b[H\x1b[2J")
	ack := true
	return ack
}

// ssh connection and workstation
func configSSH() {

}
