package handler

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"time"

)

func Heal_now(IDStimestamp string, stopshot chan bool, er, eve,noti, prog chan string ) {
	stopshot <- true
	NotiHandler("PROCEDURE MESSAGE: HEALING PROCESS HAS BEEN STARTED.....",noti)
	var detection Config
	filedata, _ := ioutil.ReadFile(filepath)
	err := json.Unmarshal(filedata, &detection)
	Errorhandler(err, red+"IP Attacker MESSAGE: Failed to Unmarshal config file"+reset,er)
	for {

		ip := findIP(er,noti)
		if len(ip) <= 0 {
			NotiHandler(red + "IP Attacker MESSAGE: Attacker IP Address not found" + reset,noti)
			continue
		} else {
			NotiHandler(red + "IP Attacker MESSAGE: ATTACKER IP ADDRESS FOUND\n STARTING SELF-HEALING PROCESS" + reset,noti)
			Close_FirewallRules(er,noti)
			break
		}
	}

	Done := CreateSnapshotTOcompare(IDStimestamp,er,eve,noti)
	if Done {
		NotiHandler("RSYNC MESSAGE: THE ROLLBACK IS COMPLETED SCECSSFULLY",noti)
		stopshot <- false
	}

	/*

	   thershold is :
	    - history
	    -  use diff command to look for hcanges between the snapshot and the current system state
	    - timestamp before or after

	*/

}

func IDS() {

}

func CreateSnapshotTOcompare(Snorttimestamp string, er, eve,noti chan string ) bool {

	filedata, _ := ioutil.ReadFile(filepath)
	snapshotlogs := filedata
	module := "snapshots"
	var backup Config
	err := json.Unmarshal(snapshotlogs, &backup)
	if err != nil {
		Errorhandler(err,red+"Error: can't Unmarshal snapshotlogs"+reset,er)
	}
	remote := fmt.Sprintf("%v@%v::%v", backup.Workstationinfo.SnapshotsUser, backup.Workstationinfo.IPaddr, module)

	list := exec.Command("sudo", "rsync", "-aAXv", "--list-only", `"--exclude={"/etc/ASRS_agent/*", "/dev/*","/proc/*","/sys/*","/tmp/*","/run/*","/mnt/*","/media/*","/lost+found"}"`, remote, "/")
	outputlist, err := list.CombinedOutput()
	if err != nil {
		Errorhandler(err,red+"RSYNC MESSAGE: Error list files to be backed up"+reset,er)
	}

	snorttimestamp := "678687687"
	isSnortafter := checktimeSN_SP(snorttimestamp, backup.Backup.Ltimestamp,er) // true or false
	linesofdiff := checkdiffANDsenfiles(outputlist,er,eve)                          // log it
	if !isSnortafter && len(linesofdiff) <= 0 {
		return false
	} else {
		create := exec.Command("sudo", "rsync", "-aAXv", "--delete", `"--exclude={"/etc/ASRS_agent/*", "/dev/*","/proc/*","/sys/*","/tmp/*","/run/*","/mnt/*","/media/*","/lost+found"}"`, remote, "/")
		_, err := create.CombinedOutput()
		if err != nil {
			Errorhandler(err,red + "RSYNC MESSAGE: Error restore backup" + reset, er)
		}
		return true
	}

	//conf.ASRS_CONF = append(conf.ASRS_CONF, )

	/*
		0 - check the last snapshots timestamp and the timestamp of the attack
		1 - make a txt file that contains the sensitive config files
		2 - load the files and put it with diff
		3 - rotate into files that changed and log it
		4 - if files changed && if history contains mal commands then rollback
		 or
		 use aide tool to check sensitive files changes and compare the change if after or before the attack
		 if after the attack then rollback -+

	*/
	// whoami
	// uname -a
	// cat /etc/passwd
	// uname -s
	// ncat -vvlp
	//
	/*
			The example above shows the sample2rs.txt file is missing at the destination.

		Possible letters in the output are:

		    f – stands for file
		    d – shows the destination file is in question
		    t – shows the timestamp has changed
		    s – shows the size has changed*/
	//this show diff between source and dest
	// rsync -avi /home/test/Desktop/Linux/ /home/test/Desktop/rsync
	// OR use --list-only to list the files that will be transfered

	// change it

}

func checkdiffANDsenfiles(diff []byte, er, eve chan string ) []string {
	senfiles := "/etc/ASRS_agent/.config/senfiles.txt"
	leno := []string{}
	conns := []string{}
	senfiles_lines := senloadlines(senfiles,er)
	//diff_lines := diffloadlines(string(diff))
	// Scan each line of the output
	scanner := bufio.NewScanner(strings.NewReader(string(diff)))
	EventHandler("STARTING CHECKING FOR BREACHED FILES",eve)
	for scanner.Scan() {
		line := scanner.Text()
		// Check if the line contains sensitive files
		for i := 0; i < len(senfiles_lines); i++ {
			if strings.Contains(line, senfiles_lines[i]) {
				leno[i] = fmt.Sprintf("Line contains %v  :%v\n", senfiles_lines[i], strings.TrimSpace(line))
				conns = append(conns, line)
				//fmt.Printf("%v\n",conns)
			}
		}
	}
	f := fmt.Sprintf("LIST OF BREACHED FILES\n%v",leno)
	EventHandler(f,eve)
	if err := scanner.Err(); err != nil {
		Errorhandler(err,"Error scanning output:",er)
	}
	return conns
}

func senloadlines(file string,er chan string  ) []string {
	filedata, err := os.Open(file)
	if err != nil {
		Errorhandler(err,"Error opening file:",er)
	}
	defer filedata.Close()

	lines := []string{}

	scanner := bufio.NewScanner(filedata)
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}

	if err := scanner.Err(); err != nil {
		Errorhandler(err,"Error scanning file:",er)
	}

	return lines
}

func checktimeSN_SP(Snort_time, snap_time string,er chan string ) bool {

	currentYear := time.Now().Year()
	Snort_time = fmt.Sprintf("%d/%s", currentYear, Snort_time)
	layout1 := "2006/01/02-15:04:05.000000"
	SNORTparsed, err := time.Parse(layout1, Snort_time)
	if err != nil {
		Errorhandler(err,"Error parsing timestamp1:", er)

	}

	layout2 := "2006-01-02 15:04:05"
	SNAPparsed, err := time.Parse(layout2, snap_time)
	if err != nil {
		Errorhandler(err,"Error parsing timestamp2:", er)

	}

	// Compare the parsed timestamps
	if SNORTparsed.Before(SNAPparsed) {
		//EventHandler(Snort_time, "is before", snap_time)
		return false
	} else if SNORTparsed.After(SNAPparsed) {
		//fmt.Println(Snort_time, "is after", snap_time)
		return true
	} else {
		//fmt.Println(Snort_time, "is equal to", snap_time)
		return true

	}
}

func Get_Status(er, eve, prog chan string ) {
	ProgHandler("PROCEDURE MESSAGE: STATUS HAS BEEN SENT",prog)

}

func Restore_Backup(done chan bool, er,noti chan string ) {
	NotiHandler("SYNC WEB FILES STARTED",noti)
	// for loop and wait for sync file and log it

	var conf Config
	databaseMOD := "database"
	websiteMOD := "website"
	filedata, _ := ioutil.ReadFile(filepath)
	err := json.Unmarshal(filedata, &conf)
	Errorhandler(err, red+"RESTORE WEB FILES MESSAGE: Failed to unmarshal config"+reset, er)
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
					cmd := exec.Command("sudo", "rsync", "-av", "--delete", remote, back)
					_, err := cmd.CombinedOutput()
					Errorhandler(err, red+"RSYNC RESTORE MESSAGE: Faild to restore webiste files to remote directory"+reset,er)
					//fmt.Println(string(outpit))
				}
			} else if i == 1 {
				for _, back := range database {
					remote := fmt.Sprintf("%v@%v::%v", conf.Workstationinfo.Webuser, conf.Workstationinfo.IPaddr, databaseMOD)
					cmd := exec.Command("sudo", "rsync", "-av", "--delete", remote, back)
					_, err := cmd.CombinedOutput()
					Errorhandler(err, red+"RSYNC RESTORE MESSAGE: Faild to restore database files to remote directory"+reset,er)
					//fmt.Println(string(outpit))
				}
			}

		}
	}()

	NotiHandler("RETORE BACKUP MESSAGE: Restoring Website Backup DONE",noti)
	done <- true
}

func Close_FirewallRules(er, noti chan string ) {
	NotiHandler(green + "FIREWALL MESSAGE: STARTING CLOSING CONNECTIONS " + reset,noti)
	var conf Config
	filedata, _ := ioutil.ReadFile(filepath)
	err := json.Unmarshal(filedata, &conf)
	Errorhandler(err, red+"FIND IP MESSAGE: Failed to unmarshal config"+reset, er)
	// allow
	_, _ = exec.Command("iptables", "-A", "INPUT", "-p", "tcp", "-s", conf.Workstationinfo.IPaddr, "-j", "ACCEPT").Output()
	_, _ = exec.Command("iptables", "-A", "OUTPUT", "-p", "tcp", "-d", conf.Workstationinfo.IPaddr, "-j", "ACCEPT").Output()
	// reject others
	_, _ = exec.Command("iptables", "-A", "INPUT", "-j", "REJECT").Output()
	_, _ = exec.Command("iptables", "-A", "OUTPUT", "-j", "REJECT").Output()
	_, _ = exec.Command("systemctl", "restart", "iptables").Output()
	NotiHandler(green + "FIREWALL MESSAGE: CLOSING CONNECTIONS COMPLETED" + reset,noti)
}

func Open_FirewallRules(er, noti chan string ) {
	NotiHandler(green + "FIREWALL MESSAGE: OPENING CONNECTIONS " + reset,noti)

	var conf Config
	filedata, _ := ioutil.ReadFile(filepath)
	err := json.Unmarshal(filedata, &conf)
	Errorhandler(err, red+"FIND IP MESSAGE: Failed to unmarshal config"+reset,er)

	// remove
	_, _ = exec.Command("iptables", "-D", "INPUT", "-p", "tcp", "-s", conf.Workstationinfo.IPaddr, "-j", "ACCEPT").Output()
	_, _ = exec.Command("iptables", "-D", "OUTPUT", "-p", "tcp", "-d", conf.Workstationinfo.IPaddr, "-j", "ACCEPT").Output()

	// relove reject
	_, _ = exec.Command("iptables", "-D", "INPUT", "-j", "REJECT").Output()
	_, _ = exec.Command("iptables", "-D", "OUTPUT", "-j", "REJECT").Output()
	_, _ = exec.Command("systemctl", "restart", "iptables").Output()
	NotiHandler(green + "FIREWALL MESSAGE: OPEINING CONNECTIONS COMPLETED" + reset,noti)

}

func findIP(er, noti chan string) []string {
	NotiHandler(green + "SEARCHING FOR ATTACKER IP CONNECTIONS" + reset,noti)

	var conf Config
	filedata, _ := ioutil.ReadFile(filepath)
	err := json.Unmarshal(filedata, &conf)
	Errorhandler(err, red+"FIND IP MESSAGE: Failed to unmarshal config"+reset,er)

	cmd := exec.Command("sh", "-c", "ss -antp")
	output, err := cmd.Output()
	if err != nil {
		Errorhandler(err,"Error executing command:",er)

	}
	conns := []string{}
	substrings := []string{"ESTAB", "LISTEN", "TIME-WAIT"}
	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	for scanner.Scan() {
		line := scanner.Text()

		for i := 0; i < 3; i++ {
			if strings.Contains(line, substrings[i]) && strings.Contains(line, conf.Detectionmarker.AttackerIP) {
				//fmt.Printf("Line contains both %v and 192.168.43.37 :%v\n", substrings[i], strings.TrimSpace(line))
				conns = append(conns, line)

				// log lines of connecitons
			}
		}
	}
	//fmt.Println(conns)
	if err := scanner.Err(); err != nil {
		Errorhandler(err,"Error scanning output:",er)
	}
	NotiHandler(green + "SEARCHING FOR ATTACKER IP CONNECTIONS COMPTELED" + reset,er)

	return conns
}
