package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os/exec"
)

func Heal_now() {
	fmt.Println("PROCEDURE MESSAGE: HEALING PROCESS HAS BEEN STARTED.....")
	var detection Config
	filedata, _ := ioutil.ReadFile(filepath)
	err := json.Unmarshal(filedata,&detection)
	errorhandler(err,red+"IP Attacker MESSAGE: Failed to Unmarshal config file"+reset)
	for {
		// netstat -antp
		//Use the netstat command to list all active connections to the Linux machine: 
	cmd := exec.Command("grep", "-o", detection.Detectionmarker.AttackerIP, "/var/log/syslog.1", "|", "head", "-n", "1")
	ip, err := cmd.Output()
	errorhandler(err, red+"IP Attacker MESSAGE: Failed to execute the command"+reset)
	if len(ip) <= 0 {
		fmt.Println(red+"IP Attacker MESSAGE: Attacker IP Address not found"+reset)
		continue
	}else {
		fmt.Println(red+"IP Attacker MESSAGE: ATTACKER IP ADDRESS FOUND\n STARTING SELF-HEALING PROCESS"+reset)
		break
	}
	}

}

func Get_Status() {
	fmt.Println("PROCEDURE MESSAGE: STATUS HAS BEEN SENT")

}

func Restore_Backup(done chan bool) {
	<-done
	fmt.Println("SYNC WEB FILES STARTED")
	// for loop and wait for sync file and log it
	var conf Config
	databaseMOD := "database"
	websiteMOD := "website"
	filedata, _ := ioutil.ReadFile(filepath)
	err := json.Unmarshal(filedata, &conf)
	errorhandler(err, red+"RESTORE WEB FILES MESSAGE: Failed to unmarshal config"+reset)
	var website = []string{
		"/var/www/html/",
		"/usr/share/nginx/html/",
	}
	var database = []string{
		"/var/lib/mysql/",
		"/var/lib/pgsql/",
	}

	//pass := "--password-file=/etc/ASRS_agent/.config/pass.txt"

	go func() {
		for i := 0; i < 2; i++ {

			if i == 0 {
				for _, back := range website {
					remote := fmt.Sprintf("%v@%v::%v", conf.Workstationinfo.Webuser, conf.Workstationinfo.IPaddr, websiteMOD)
					cmd := exec.Command("sudo", "rsync", "-av", "--delete", remote,back)
					outpit, err := cmd.CombinedOutput()
					errorhandler(err, red+"RSYNC RESTORE MESSAGE: Faild to restore webiste files to remote directory"+reset)
					fmt.Println(string(outpit))
				}
			} else if i == 1 {
				for _, back := range database {
					remote := fmt.Sprintf("%v@%v::%v", conf.Workstationinfo.Webuser, conf.Workstationinfo.IPaddr, databaseMOD)
					cmd := exec.Command("sudo", "rsync", "-av", "--delete",remote, back)
					outpit, err := cmd.CombinedOutput()
					errorhandler(err, red+"RSYNC RESTORE MESSAGE: Faild to restore database files to remote directory"+reset)
					fmt.Println(string(outpit))
				}
			}

		}
	}()
	fmt.Println("restore backup done")
	done <- true
}
