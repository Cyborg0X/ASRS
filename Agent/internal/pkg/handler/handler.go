package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"os/exec"
	"sync"
	"time"
)

var red = "\033[31m"
var green = "\033[32m"
var reset = "\033[0m"

var k int
var filepath = "/etc/ASRS_agent/.config/config.json"

type DataType int

const (
	TypeA1 DataType = iota
	TypeA2
)

type A1 struct {
	A string `json:"procedure"`
}

type A2 struct {
	A          string `json:"procedure"`
	AttackerIP string `json:"Attacker IP"`
	AttackTime string `json:"Time of attack"`
}

/*
	type SSH struct {
		Proceduree string `json:"procedure"`
		Username   string `json:"SSH username"`
		Pass       string `json:"SSH pass"`
	}
*/
type DataWrapper struct {
	Type DataType    `json:"type"`
	Data interface{} `json:"data"`
}

func TaskHandler(wg *sync.WaitGroup, chanconn chan net.Conn) {
	EventHandler("ASRS TASK HANDLER STARTED")
	stopshot := make(chan bool, 1)
	defer wg.Done()
	B3 := DetectionMarker()
	if B3 {
		wg.Add(1)
		backchan := make(chan bool)
		EventHandler("ASRS RESTORE BACKUP STARTED")
		go Restore_Backup(backchan)
		<-backchan
		ProcedureHandler(wg, chanconn, B3, stopshot)
		var detected Config
		filedata, _ := ioutil.ReadFile(filepath)
		err := json.Unmarshal(filedata, &detected)
		Errorhandler(err, red+"RESTORE BACKUP MESSAGE: Failed to unmarshal marker:"+reset)
		detected.Detectionmarker.Markerisdetected = false
		json.MarshalIndent(filedata, "", "  ")
	}

	wg.Add(2)

	go Local_actions(wg, stopshot)
	go ProcedureHandler(wg, chanconn, B3, stopshot)
	wg.Wait()

}

func Response_Sender(message string, conn net.Conn) {
	EventHandler("RESPONSE SENDER IN ACTION")
	for {
		_, err := conn.Write([]byte(message))
		if err != nil {
			Errorhandler(err, "SENDER MESSAGE: Faild to send")
			continue
		}
		conn.Close()
		break
	}
}

func ProcedureHandler(wg *sync.WaitGroup, chanconn chan net.Conn, B3 bool, stopshot chan bool) {
	EventHandler("PROCEDURE HANDLER STARTED")
	fmt.Println()
	defer wg.Done()

	for {

		if conneceted, ok := <-chanconn; ok {

			message := make([]byte, 1024)
			n, err := conneceted.Read(message)
			if err != nil {
				Errorhandler(err, fmt.Sprintf("Failed to read from connection: %v", err))
				break
			}
			defer conneceted.Close()
			var wrapper DataWrapper
			k = k + 1
			err = json.Unmarshal(message[:n], &wrapper)
			Errorhandler(err, red+"Can't unmarshal received message"+reset)
			//dataStr := fmt.Sprintf("%v", wrapper.Data)
			//fmt.Println(wrapper.Data)
			switch wrapper.Type {
			case TypeA1:
				dataMap := wrapper.Data.(map[string]interface{})
				if dataMap["procedure"] == "A1" {
					ProgHandler("PROCEDURE A1 RECEIVED")
					if B3 {
						Response_Sender("B3PROC", conneceted)
						B3 = false
						return
					}
					go Get_Status()

				}
			case TypeA2:
				dataMap := wrapper.Data.(map[string]interface{})
				ProgHandler("PROCEDURE A2 RECEIVED")
				attacker := dataMap["Attacker IP"]
				time := dataMap["Time of attack"]
				AttackerIP(attacker.(string), time.(string))
				// mar
				if B3 {
					Response_Sender("B3PROC", conneceted)
					B3 = false
					return
				}
				go Heal_now(time.(string), stopshot)

			}

		}
	}

}

func AttackerIP(ip string, time string) {

	var marsh Config
	filedata, err := ioutil.ReadFile("/etc/ASRS_agent/.config/config.json")
	if err != nil {
		fmt.Println(red+"Error:"+reset, err)
		return
	}
	err = json.Unmarshal(filedata, &marsh)
	if err != nil {
		fmt.Println(red+"Error:"+reset, err)
		return
	}
	marsh.Detectionmarker.AttackerIP = ip
	marsh.Detectionmarker.AttackTiming = time
	conf, err := json.MarshalIndent(marsh, "", "  ")
	Errorhandler(err, red+"Attacker IP MESSAGE: Can't Marshal IP ADDRESS and Time of the attack"+reset)
	fileper, _ := os.Stat(filepath)
	per := fileper.Mode().Perm()
	_ = ioutil.WriteFile("/etc/ASRS_agent/.config/config.json", conf, per)
	fmt.Println(green + "Attacker IP MESSAGE: Attacker IP saved" + reset)

}

func Local_actions(wg *sync.WaitGroup, stopshot chan bool) {
	EventHandler("LOCAL ACTIONS STARTED")
	// receive channel from B2 to terminate
	wg.Add(1)
	defer wg.Done()

	go func() {
		cg := make(chan bool, 1)
		for {
			go CreateSnapshot(cg, stopshot)
			<-cg
		}
	}()
	go func() {
		for {
			time.Sleep(time.Second * 20)
			ProgHandler("SYNCING WEB FILES...")
			//Sync_web_files()
		}

	}()
	wg.Wait()

}

func CreateSnapshot(vx chan bool, stopshot chan bool) {
	ProgHandler("ASRS CREATING SNAPSHOT...")
	fmt.Println()
	var counter int
	//asrs_conf := "ASRS_CONF"
	//discription := "Incremental Backup"
	mountpoint := "/home/agent/"
	module := "snapshots"

	filedata, _ := ioutil.ReadFile(filepath)
	snapshotlogs := filedata

	var checker Config
	err := json.Unmarshal(snapshotlogs, &checker)
	if err != nil {
		Errorhandler(err, "BACKUP: Can't Unmarshal snapshot logs")
		return
	}
	remote := fmt.Sprintf("%v@%v::%v", checker.Workstationinfo.SnapshotsUser, checker.Workstationinfo.IPaddr, module)

	if !checker.Backup.FullSnapshot {
		EventHandler("RSYNC Started taking full backup for your system files, this process may take time please wait.....")
		config := exec.Command("sudo", "rsync", "-aAXv", "--delete", mountpoint, `"--exclude={"/etc/ASRS_agent/*", "/dev/*","/proc/*","/sys/*","/tmp/*","/run/*","/mnt/*","/media/*","/lost+found"}"`, remote)
		output, err := config.CombinedOutput()
		if err != nil {
			Errorhandler(err, "RSYNC MESSAGE: Failed to take first backup:")
			vx <- true
			return
		}
		fmt.Println(string(output)) // log it
		now := time.Now()
		checker.Backup.Ltimestamp = now.Format("2006-01-02 15:04:05")
		checker.Backup.FullSnapshot = true
		done, err := json.MarshalIndent(checker, "", "  ")
		if err != nil {
			Errorhandler(err, "RSYNC MESSAGE: Failed to write to snapshot checker")
			vx <- true
			return
		}
		fileper, _ := os.Stat(filepath)
		per := fileper.Mode().Perm()
		ioutil.WriteFile(filepath, done, per)
		ProgHandler("RSYNC MESSAGE: Backup completed ....." )

	}

	for {
		time.Sleep(time.Minute * 4)
		select {
		case value := <-stopshot:
			if value {
				<-stopshot
			}
		default:
			ProgHandler("RSYNC Started taking backup...")
		}
		checker.Detectionmarker.Markerisdetected = true
		Updated_Marker, err := json.MarshalIndent(checker, "", "  ")
		if err != nil {
			now := time.Now()
			Errorhandler(err,fmt.Sprintf("RSYNC MESSAGE: Failed to write detection marker to `true`\n detection marker is false in the backup the has been taken in this time %v\n", now.Format("2006-01-02 15:04:05")))
			vx <- true
			return
		}
		fileper, _ := os.Stat(filepath)
		per := fileper.Mode().Perm()
		err = ioutil.WriteFile(filepath, Updated_Marker, per)
		counter++
		time.Sleep(time.Second * 1)

		create := exec.Command("sudo", "rsync", "-aAXv", "--delete", mountpoint, `"--exclude={"/etc/ASRS_agent/*", "/dev/*","/proc/*","/sys/*","/tmp/*","/run/*","/mnt/*","/media/*","/lost+found"}"`, remote)
		output, err := create.CombinedOutput()
		if err != nil {
			Errorhandler(err,fmt.Sprintf(red+"RSYNC Error creating backup number: %v\n ERROR: %v \n output: %v\n"+reset, counter, err, string(output)))
			
			vx <- true
			return
		}
		now := time.Now()
		checker.Backup.Ltimestamp = now.Format("2006-01-02 15:04:05")
		checker.Detectionmarker.Markerisdetected = false
		done, err := json.MarshalIndent(checker, "", "  ")
		if err != nil {
			Errorhandler(err,"SNAPPER MESSAGE: Failed to write detection marker to false")
			vx <- true
			return
		}

		ioutil.WriteFile(filepath, done, per)
		ProgHandler("RSYNC MESSAGE: SYSTEM FILES SYNCED")

		//remotepath := "/etc/ASRS_WS/.database/snapshots_backup/"
		//fmt.Println(remote)
		//fmt.Println(string(output)) // log it later ALSO set JSON OUTPUT FORMAT IN SNAPPER
		//pass := "--password-file=/etc/ASRS_agent/.config/pass.txt"

	}
}

//for loop, wait for 1 hour, set detection marker, take snapshot, remove detection marker
// to list snapshots of config >>> sudo snapper -c ASRS_CONF list

// to set new default config >>>

func DetectionMarker() bool {
	var detector Config
	if detector.Detectionmarker.Markerisdetected {
		EventHandler("DETECTION MARKER MESSAGE: DETELCTION MARKER EXIST")
		return detector.Detectionmarker.Markerisdetected

	}
	return false
}

func ProcedureReceiver() {

}

/*
func Sync_web_files() {
	fmt.Println("SYNC WEB FILES STARTED")
	// for loop and wait for sync file and log it
	var conf Config
	databaseMOD := "database"
	websiteMOD := "website"
	filedata, _ := ioutil.ReadFile(filepath)
	err := json.Unmarshal(filedata, &conf)
	Errorhandler(err, red+"SYNC WEB FILES MESSAGE: Failed to unmarshal config"+reset)
	var website = []string{
		"/var/www/html/",
		"/usr/share/nginx/html/",
	}
	var database = []string{
		"/var/lib/mysql/",
		"/var/lib/pgsql/",
	}
	//pass := "--password-file=/etc/ASRS_agent/.config/pass.txt"

	func() {
		for i := 0; i < 2; i++ {

			if i == 0 {
				for _, back := range website {
					remote := fmt.Sprintf("%v@%v::%v", conf.Workstationinfo.Webuser, conf.Workstationinfo.IPaddr, websiteMOD)
					cmd := exec.Command("sudo", "rsync", "-av", "--delete", back, remote)
					outpit, err := cmd.CombinedOutput()
					Errorhandler(err, red+"RSYNC MESSAGE: Faild to sync webiste files to remote directory"+reset)
					fmt.Println(string(outpit))
				}
			} else if i == 1 {
				for _, back := range database {
					remote := fmt.Sprintf("%v@%v::%v", conf.Workstationinfo.Webuser, conf.Workstationinfo.IPaddr, databaseMOD)
					cmd := exec.Command("sudo", "rsync", "-av", "--delete", back, remote)
					outpit, err := cmd.CombinedOutput()
					Errorhandler(err, red+"RSYNC MESSAGE: Faild to sync database files to remote directory"+reset)
					fmt.Println(string(outpit))
				}
			}

		}
	}()

}

	func get_username(username string, pass string) {
		fmt.Println("GET SSH USERNAME STARTED")
		var put Config
		passw := "/etc/ASRS_agent/.config/pass.txt"
		var filesdata, _ = ioutil.ReadFile(filepath)
		file := filesdata
		_ = json.Unmarshal(file, &put)
		put.Workstationinfo.SSH_username = username
		put.Workstationinfo.SSHpass = username
		jsondata, err := json.MarshalIndent(put, "", "  ")
		Errorhandler(err, red+"can't marshal username"+reset)
		ioutil.WriteFile(filepath, jsondata, 0766)
		ioutil.WriteFile(passw, []byte(pass), 0766)
	}

func SSHkeys() {

}
*/
