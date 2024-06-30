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

var k int

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
		go Local_actions()
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
					fmt.Println("Failed to read from connection:", err)
					break
				}
				defer conneceted.Close()

				var wrapper DataWrapper
				k = k + 1

				err = json.Unmarshal(message[:n], &wrapper)
				errorhandler(err, "can't unmarshal received message")
				//dataStr := fmt.Sprintf("%v", wrapper.Data)
				//fmt.Println(wrapper.Data)
				switch wrapper.Type {
				case Typemessage:
					dataMap := wrapper.Data.(map[string]interface{})
					if dataMap["procedure"] == "A1" {
						fmt.Println("PROCEDURE MESSAGE: A1 RECEIVED")
						go Get_Status()

					} else if dataMap["procedure"] == "A2" {
						fmt.Println("PROCEDURE MESSAGE: A2 RECEIVED")
						go Heal_now()
					}
				case TypeSSH:
					dataMap := wrapper.Data.(map[string]interface{})
					userbame := dataMap["SSH username"].(string)

					go get_username(userbame)
					fmt.Println("SSH MESSAGE: SSH username RECEIVED")

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
	filepath := "/etc/ASRS_agent/.config/config.json"
	snapshotlogs, err := ioutil.ReadFile(filepath)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	var checker Config
	err = json.Unmarshal(snapshotlogs, &checker)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	if !checker.Backup.FullSnapshot {
		config := exec.Command("sudo", "snapper", "-c", asrs_conf, "create-config", mountpoint)
		output, err := config.CombinedOutput()
		if err != nil {
			fmt.Printf("Failed to initiate first snapshot: %s\n", err)
			fmt.Println(string(output))
			return
		}
		fmt.Println(string(output)) // log it

		checker.Backup.FullSnapshot = true
		done, err := json.MarshalIndent(checker, "", "  ")
		if err != nil {
			fmt.Println("Failed to write to snapshot checker")
		}
		ioutil.WriteFile(filepath, done, 0766)

		FULL := exec.Command("sudo", "snapper", "-c", asrs_conf, "create", "-t", "pre", "-d", "BASELINE SNAPSHOT")
		output_full, err := FULL.CombinedOutput()
		if err != nil {
			fmt.Println("Failed to get output of pre snaphsot creation")
		}
		fmt.Println(string(output_full))

	}
	for {
		checker.Detectionmarker.Markerisdetected = true
		Updated_Marker, err := json.MarshalIndent(checker, "", "  ")
		if err != nil {
			date := time.Now()
			fmt.Printf("Failed to write detection marker to `true`\n detection marker is false in the snapshot the has been taken in this time %v\n", date)
		}
		err = ioutil.WriteFile(filepath, Updated_Marker, 0766)
		counter++
		create := exec.Command("sudo", "snapper", "-c", asrs_conf, "create", "-d", discription, "--output", "json")
		output, err := create.CombinedOutput()
		if err != nil {
			fmt.Printf("Error creating snapshot number: %v\n ERROR: %v \n output: %v\n", counter, err, string(output))
		}

		checker.Backup.SnapshotNum = checker.Backup.SnapshotNum + 1
		checker.Detectionmarker.Markerisdetected = false
		done, err := json.MarshalIndent(checker, "", "  ")
		if err != nil {
			fmt.Println("Failed to write new snapshot number")
		}
		ioutil.WriteFile(filepath, done, 0766)

		fmt.Println(string(output)) // log it later ALSO set JSON OUTPUT FORMAT IN SNAPPER
		rsynco := exec.Command("sudo", "rsync", "-aAXv", "--delete", "/var/lib/snapper/configs/root", "/var/lib/snapper/snapshots/root/", "username@192.168.1.100:/path/to/remote/snapper/")
		routput, err := rsynco.Output()
		fmt.Println(string(routput))
		time.Sleep(time.Hour)

	}
	//for loop, wait for 1 hour, set detection marker, take snapshot, remove detection marker

}

func get_username(username string) {

	filepath := "/etc/ASRS_agent/.config/config.json"
	var put Config
	//
	file, err := ioutil.ReadFile(filepath)
	_ = json.Unmarshal(file, &put)
	put.Workstationinfo.SSH_username = username
	jsondata, err := json.MarshalIndent(put, "", "  ")
	errorhandler(err, "can't narshal username")
	ioutil.WriteFile(filepath, jsondata, 0766)
}

func ProcedureReceiver() {

}

func Local_actions() {
	// receive channel from B2 to terminate
	go func() {
		for {
			CreateSnapshot()
			time.Sleep(time.Hour)
		}
	}()
	go func() {
		Sync_web_files()
		time.Sleep(time.Minute)
	}()

}

func Sync_web_files() {
	// for loop and wait for sync file and log it
}
