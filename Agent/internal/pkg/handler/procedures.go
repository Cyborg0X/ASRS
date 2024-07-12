package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os/exec"
)

func Heal_now() {
	fmt.Println("PROCEDURE MESSAGE: HEALING PROCESS HAS BEEN STARTED.....")
	

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
