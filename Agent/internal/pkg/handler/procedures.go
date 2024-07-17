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

func Heal_now(IDStimestamp string) {
	fmt.Println("PROCEDURE MESSAGE: HEALING PROCESS HAS BEEN STARTED.....")
	var detection Config
	filedata, _ := ioutil.ReadFile(filepath)
	err := json.Unmarshal(filedata, &detection)
	errorhandler(err, red+"IP Attacker MESSAGE: Failed to Unmarshal config file"+reset)
	for {

		ip := findIP()
		if len(ip) <= 0 {
			fmt.Println(red + "IP Attacker MESSAGE: Attacker IP Address not found" + reset)
			continue
		} else {
			fmt.Println(red + "IP Attacker MESSAGE: ATTACKER IP ADDRESS FOUND\n STARTING SELF-HEALING PROCESS" + reset)
			Close_FirewallRules()
			break
		}
	}
	

/*

thershold is :
 - history 
 -  use diff command to look for hcanges between the snapshot and the current system state 
 - timestamp before or after 

*/


	
	

}

func IDS()  {
	
}

func CreateSnapshotTOcompare(Snorttimestamp string) bool {

	filedata, _ := ioutil.ReadFile(filepath)
	snapshotlogs := filedata
	module := "snapshots"
	var backup Config
	err := json.Unmarshal(snapshotlogs, &backup)
	if err != nil {
		fmt.Println(red+"Error: can't Unmarshal snapshotlogs"+reset, err)
	}
	remote := fmt.Sprintf("%v@%v::%v", backup.Workstationinfo.SnapshotsUser, backup.Workstationinfo.IPaddr, module)

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

	list := exec.Command("sudo", "rsync", "-aAXv","--list-only", `"--exclude={"/dev/*","/proc/*","/sys/*","/tmp/*","/run/*","/mnt/*","/media/*","/lost+found"}"`, remote, "/" )
	outputlist, err := list.CombinedOutput()
	if err != nil {
		fmt.Printf(red+"RSYNC MESSAGE: Error list files to be backed up, ERROR: %v \n output: %v\n"+reset, err, string(outputlist))
	}

	
	
	
	snorttimestamp := "678687687"
	isSnortafter := checktimeSN_SP(snorttimestamp, backup.Backup.Ltimestamp) // true or false
	linesofdiff := checkdiffANDsenfiles(outputlist)
	if !isSnortafter && len(linesofdiff) <= 0 {
		return false
	}else {
		create := exec.Command("sudo", "rsync", "-aAXv","--delete",`"--exclude={"/dev/*","/proc/*","/sys/*","/tmp/*","/run/*","/mnt/*","/media/*","/lost+found"}"`, remote, "/" )
		output, err := create.CombinedOutput()
		if err != nil {
			fmt.Printf(red+"RSYNC MESSAGE: Error restore backup, ERROR: %v \n output: %v\n"+reset, err, string(output))
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
	
	
}

func checkdiffANDsenfiles(diff []byte) ([]string) {
	senfiles :="/etc/ASRS_agent/.config/senfiles.txt"

	conns := []string{}
	senfiles_lines := senloadlines(senfiles)
	//diff_lines := diffloadlines(string(diff))
	// Scan each line of the output
	scanner := bufio.NewScanner(strings.NewReader(string(diff)))
	for scanner.Scan() {
		line := scanner.Text()
		// Check if the line contains sensitive files
		for i := 0; i < len(senfiles_lines); i++ {
			if strings.Contains(line, senfiles_lines[i]) {
				fmt.Printf("Line contains %v  :%v\n", senfiles_lines[i], strings.TrimSpace(line))
				conns = append(conns, line)
				//fmt.Printf("%v\n",conns)
			}
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Error scanning output:", err)
	}
	return conns
}

func senloadlines(file string) ([]string) {
	filedata, err := os.Open(file)
	if err != nil {
		fmt.Println("Error opening file:", err)
	}
	defer filedata.Close()

	lines := []string{}

	scanner := bufio.NewScanner(filedata)
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error scanning file:", err)
	}

	return lines
}




func checktimeSN_SP(Snort_time,snap_time string)(bool) {

	currentYear := time.Now().Year()
	Snort_time = fmt.Sprintf("%d/%s", currentYear, Snort_time)
	layout1 :="2006/01/02-15:04:05.000000"
	SNORTparsed, err := time.Parse(layout1, Snort_time)
	if err != nil {
		fmt.Println("Error parsing timestamp1:", err)
		
	}

	layout2 := "2006-01-02 15:04:05"
	SNAPparsed, err := time.Parse(layout2, snap_time)
	if err != nil {
		fmt.Println("Error parsing timestamp2:", err)
		
	}

	// Compare the parsed timestamps
	if SNORTparsed.Before(SNAPparsed) {
		fmt.Println(Snort_time, "is before", snap_time)
		return false
	} else if SNORTparsed.After(SNAPparsed) {
		fmt.Println(Snort_time, "is after", snap_time)
		return true
	} else {
		fmt.Println(Snort_time, "is equal to", snap_time)
		return true

	}
}


func Get_Status() {
	fmt.Println("PROCEDURE MESSAGE: STATUS HAS BEEN SENT")

}



/*	
func Restore_Backup(done chan bool) {
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
					cmd := exec.Command("sudo", "rsync", "-av", "--delete", remote, back)
					outpit, err := cmd.CombinedOutput()
					errorhandler(err, red+"RSYNC RESTORE MESSAGE: Faild to restore webiste files to remote directory"+reset)
					fmt.Println(string(outpit))
				}
			} else if i == 1 {
				for _, back := range database {
					remote := fmt.Sprintf("%v@%v::%v", conf.Workstationinfo.Webuser, conf.Workstationinfo.IPaddr, databaseMOD)
					cmd := exec.Command("sudo", "rsync", "-av", "--delete", remote, back)
					outpit, err := cmd.CombinedOutput()
					errorhandler(err, red+"RSYNC RESTORE MESSAGE: Faild to restore database files to remote directory"+reset)
					fmt.Println(string(outpit))
				}
			}

		}
	}()

	fmt.Println("RETORE BACKUP MESSAGE: Restoring Website Backup DONE")
	done <- true
}
*/


func Close_FirewallRules() {
	// allow
	_, _ = exec.Command("iptables", "-A", "INPUT", "-p", "tcp", "-s", "140.82.112.25", "-j", "ACCEPT").Output()
	_, _ = exec.Command("iptables", "-A", "OUTPUT", "-p", "tcp", "-d", "140.82.112.25", "-j", "ACCEPT").Output()
	// reject others
	_, _ = exec.Command("iptables", "-A", "INPUT", "-j", "REJECT").Output()
	_, _ = exec.Command("iptables", "-A", "OUTPUT", "-j", "REJECT").Output()
	_, _ = exec.Command("systemctl", "restart", "iptables").Output()
}

func Open_FirewallRules() {
	// remove
	_, _ = exec.Command("iptables", "-D", "INPUT", "-p", "tcp", "-s", "140.82.112.25", "-j", "ACCEPT").Output()
	_, _ = exec.Command("iptables", "-D", "OUTPUT", "-p", "tcp", "-d", "140.82.112.25", "-j", "ACCEPT").Output()

	// relove reject
	_, _ = exec.Command("iptables", "-D", "INPUT", "-j", "REJECT").Output()
	_, _ = exec.Command("iptables", "-D", "OUTPUT", "-j", "REJECT").Output()
	_, _ = exec.Command("systemctl", "restart", "iptables").Output()
}

func findIP() []string {

	var conf Config
	filedata, _ := ioutil.ReadFile(filepath)
	err := json.Unmarshal(filedata, &conf)
	errorhandler(err, red+"FIND IP MESSAGE: Failed to unmarshal config"+reset)


	cmd := exec.Command("sh", "-c", "ss -antp")
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Error executing command:", err)

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
	fmt.Println(conns)
	if err := scanner.Err(); err != nil {
		fmt.Println("Error scanning output:", err)
	}
	return conns
}






