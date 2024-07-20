package checker

import (
	//"encoding/json"
	"fmt"
	"os"

	//"io/ioutil"
	"os/exec"
	"strings"

	//"github.com/Cyborg0X/ASRS/Agent/internal/pkg/communication"
	"github.com/Cyborg0X/ASRS/Agent/internal/pkg/handler"
)

var red = "\033[31m"
var green = "\033[32m"
var reset = "\033[0m"

func Depcheck() bool {
	// return ack of
	packages := make([]string, 1)
	packages[0] = "rsync"

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
			fmt.Printf("installing %v..... please wait\n", packages[i])
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
			fmt.Printf(green+"%v is installed\n"+reset, key)
			feelsgood++
		} else {
			fmt.Println(red+"Some packages not intalled!!!!"+reset)
			fmt.Println(key, value)
		} 
	}

	if feelsgood == 1 {
		fmt.Println(green+"ALL PACKAGES HAS BEEN INSTALLED SECCUSSFULLY"+reset)
		fmt.Println(green+"Initilizing config file,,,,,,"+reset)
		filepath := "/etc/ASRS_agent/.config/config.json"

		
		file, _ := os.Stat(filepath)
		// If file doesn't exist, assume it's the first run and create a new one with InitializeJSON
		if  file.Size() == 0 {
			err := handler.InitializeJSON()
			if err != nil {
				fmt.Println(red+"Error initialize config file"+reset,err)
			}
		}
		var detector handler.Config
		fmt.Println(green+"marker : "+reset, detector.Detectionmarker)


	} else {
		panic(red+"Error checking startup"+reset)
	}
	//fmt.Fprint(os.Stdout, "\x1b[H\x1b[2J")
	ack := true
	return ack
}
