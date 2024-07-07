package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"os/exec"
	"sync"
	"time"
)

var red = "\033[31m"
var green = "\033[32m"
var reset = "\033[0m"

var k int
var filepath = "/etc/ASRS_agent/.config/config.json"
var filedata, _ = ioutil.ReadFile(filepath)

type DataType int

const (
	Typemessage DataType = iota
	TypeSSH
)

type A1A2 struct {
	A string `json:"procedure"`
}

type SSH struct {
	Proceduree string `json:"procedure"`
	Username   string `json:"SSH username"`
}

type DataWrapper struct {
	Type DataType    `json:"type"`
	Data interface{} `json:"data"`
}

func TaskHandler(wg *sync.WaitGroup, chanconn chan net.Conn, B3 bool) {

	defer wg.Done()
	if B3 {
		backchan := make(chan bool)
		go Restore_Backup(backchan)
		<-backchan
		return
	}

	wg.Add(1)
	if !B3 {
		go Local_actions(wg)
	}
	go ProcedureHandler(wg, chanconn)
	wg.Wait()

}

func ProcedureHandler(wg *sync.WaitGroup, chanconn chan net.Conn) {
	defer wg.Done()

	for {

		if conneceted, ok := <-chanconn; ok {

			for {
				message := make([]byte, 1024)
				n, err := conneceted.Read(message)
				if err != nil {
					fmt.Println(red + "Failed to read from connection:" + reset, err)
					break
				}
				defer conneceted.Close()

				var wrapper DataWrapper
				k = k + 1

				err = json.Unmarshal(message[:n], &wrapper)
				errorhandler(err, red + "can't unmarshal received message" + reset)
				//dataStr := fmt.Sprintf("%v", wrapper.Data)
				//fmt.Println(wrapper.Data)
				switch wrapper.Type {
				case Typemessage:
					dataMap := wrapper.Data.(map[string]interface{})
					if dataMap["procedure"] == "A1" {
						fmt.Println( green+"PROCEDURE MESSAGE: A1 RECEIVED"+reset)
						go Get_Status()

					} else if dataMap["procedure"] == "A2" {
						fmt.Println(green+"PROCEDURE MESSAGE: A2 RECEIVED"+reset)
						go Heal_now()
					}
				case TypeSSH:
					dataMap := wrapper.Data.(map[string]interface{})
					userbame := dataMap["SSH username"].(string)
					go get_username(userbame)
					fmt.Println(green+"SSH MESSAGE: SSH username RECEIVED"+reset)
					// sending keys for SSH rsync
					keys := SSH_config()
					conneceted.Write([]byte(keys))
					conneceted.Close()

				}
				break
			}

		}
	}

}

func CreateSnapshot() {
	var counter int
	asrs_conf := "ASRS_CONF"
	discription := "Incremental Backup"
	mountpoint := "/"
	snapshotlogs := filedata

	var checker Config
	err := json.Unmarshal(snapshotlogs, &checker)
	if err != nil {
		fmt.Println(red+"Error: can't Unmarshal snapshotlogs"+reset, err)
		return
	}
	if !checker.Backup.FullSnapshot {

		config := exec.Command("sudo", "snapper", "-c", asrs_conf, "create-config", mountpoint)
		output, err := config.CombinedOutput()
		if err != nil {
			fmt.Printf(red+"SNAPPER MESSAGE: Failed to initiate first snapshot: %s\n"+reset, err)
			fmt.Println(string(output))
			return
		}
		fmt.Println(string(output)) // log it

		checker.Backup.FullSnapshot = true
		done, err := json.MarshalIndent(checker, "", "  ")
		if err != nil {
			fmt.Println(red+"SNAPPER MESSAGE: Failed to write to snapshot checker"+reset)
		}
		ioutil.WriteFile(filepath, done, 0766)

		FULL := exec.Command("sudo", "snapper", "-c", asrs_conf, "create", "-t", "pre", "-d", "BASELINE SNAPSHOT")
		output_full, err := FULL.CombinedOutput()
		if err != nil {
			fmt.Println(red+"SNAPPER MESSAGE: Failed to get output of pre snaphsot creation"+reset)
		}
		fmt.Println(string(output_full))

	}
	for {
		checker.Detectionmarker.Markerisdetected = true
		Updated_Marker, err := json.MarshalIndent(checker, "", "  ")
		if err != nil {
			date := time.Now()
			fmt.Printf(red+"SNAPPER MESSAGE: Failed to write detection marker to `true`\n detection marker is false in the snapshot the has been taken in this time %v\n"+reset, date)
		}
		err = ioutil.WriteFile(filepath, Updated_Marker, 0766)
		counter++

		create := exec.Command("sudo", "snapper", "-c", asrs_conf, "create", "-d", discription)
		output, err := create.CombinedOutput()
		if err != nil {
			fmt.Printf(red+"SNAPPER MESSAGE: Error creating snapshot number: %v\n ERROR: %v \n output: %v\n"+reset, counter, err, string(output))
		}

		checker.Backup.SnapshotNum = checker.Backup.SnapshotNum + 1
		checker.Detectionmarker.Markerisdetected = false
		done, err := json.MarshalIndent(checker, "", "  ")
		if err != nil {
			fmt.Println(red+"SNAPPER MESSAGE: Failed to write new snapshot number"+reset)
		}
		ioutil.WriteFile(filepath, done, 0766)
		remotepath := "/etc/ASRS_WS/.database/snapshots_backup"
		remote := fmt.Sprintf("%v@ %v:%v", checker.Workstationinfo.SSH_username, checker.Workstationinfo.IPaddr, remotepath)
		fmt.Println(string(output)) // log it later ALSO set JSON OUTPUT FORMAT IN SNAPPER
		rsynco := exec.Command("sudo", "rsync", "-av", "--delete", "/.snapshots", remote)
		routput, err := rsynco.Output()
		errorhandler(err, red+"SNAPPER MESSAGE: Faild to sync snapshots"+reset)
		fmt.Println(string(routput))
		time.Sleep(time.Minute * 2)

	}
	//for loop, wait for 1 hour, set detection marker, take snapshot, remove detection marker
	// to list snapshots of config >>> sudo snapper -c ASRS_CONF list 
	
	// to set new default config >>>  




}

func get_username(username string) {

	var put Config
	//
	file := filedata
	_ = json.Unmarshal(file, &put)
	put.Workstationinfo.SSH_username = username
	jsondata, err := json.MarshalIndent(put, "", "  ")
	errorhandler(err, red+"can't narshal username"+reset)
	ioutil.WriteFile(filepath, jsondata, 0766)
}

func SSHkeys() {

}

func ProcedureReceiver() {

}

func Local_actions(wg *sync.WaitGroup) {
	// receive channel from B2 to terminate
	defer wg.Done()
	go func() {
		for {
			CreateSnapshot()
			time.Sleep(time.Minute * 2)
		}
	}()
	go func() {
		for {
			time.Sleep(time.Second * 5)
			Sync_web_files()
		}

	}()

}

func Sync_web_files() {
	// for loop and wait for sync file and log it
	var conf Config
	data := filedata
	err := json.Unmarshal(data, &conf)
	errorhandler(err, red+"SYNC WEB FILES MESSAGE: Failed to unmarshal config"+reset)
	var website = []string{
		"/var/www/html/",
		"/usr/share/nginx/html/",
	}
	var database = []string{
		"/var/lib/mysql/",
		"/var/lib/pgsql/",
	}

	var WSdir = []string{
		"/etc/ASRS_WS/.database/database_backup",
		"/etc/ASRS_WS/.database/website_backup",
	}
	go func() {
		for i, dir := range WSdir {

			dest := fmt.Sprintf("%v@%v:%v", conf.Workstationinfo.SSH_username, conf.Workstationinfo.IPaddr, dir)
			if i == 0 {
				for _, back := range website {

					cmd := exec.Command("sudo", "rsync", "-av", "--delete", back, dest)
					outpit, err := cmd.Output()
					errorhandler(err, red+"RSYNC MESSAGE: Faild to sync webiste files to remote directory"+reset)
					fmt.Println(string(outpit))
				}
			} else if i == 1 {
				for _, back := range database {

					cmd := exec.Command("sudo", "rsync", "-av", "--delete", back, dest)
					outpit, err := cmd.Output()
					errorhandler(err, red+"RSYNC MESSAGE: Faild to sync database files to remote directory"+reset)
					fmt.Println(string(outpit))
				}
			}

		}
	}()

}
