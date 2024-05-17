package checker

import (
	//"encoding/json"
	"fmt"
	//"io/ioutil"
	"os/exec"
	"strings"

	"github.com/Cyborg0X/ASRS/Agent/internal/pkg/communication"
	"github.com/Cyborg0X/ASRS/Agent/internal/pkg/handler"
)

func Depcheck() (ack bool) {
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
		install := exec.Command("sudo", "apt", "install", packages[i], "-y")
		output, err := checkpkg.CombinedOutput()
		outputstr := string(output)
		if err != nil {
			checklist[packages[i]] = "NOT installed"
			//fmt.Printf("installing %v..... please wait\n", packages[i])
			install.Run()

		} else if strings.Contains(outputstr, "Status: install ok installed") {

			checklist[packages[i]] = "installed"

		}
		// don't forget to add some lines to interact with a text-based user interface (TUI)
	}
	var feelsgood int
	for key, value := range checklist {
		//fmt.Printf("%v %v\n", key, value)
		if value == "installed" {
			fmt.Printf("%v is installed\n", key)
			feelsgood++
		} else {
			fmt.Println("Some packages not intalled!!!!")
			fmt.Println(key, value)
		}
	}

	if feelsgood == 6 {
		fmt.Println("ALL PACKAGES HAS BEEN INSTALLED SECCUSSFULLY")
		fmt.Println("Initilizing config file,,,,,,")
		err := handler.InitializeJSON()
		if err != nil {
			fmt.Println("Error initilizing config file", err)
		}
		communication.Checkandassign_wsip()
		
	} else {
		panic("Please try to install the package manually")
	}
	
	ack = true
	return ack

}


// ssh connection and workstation
func configSSH() {

}
