package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
)



type Config struct {
	Agentinfo struct {
		Ipaddr string `json:"AGIP"`
		Port   string `json:"AGport"`
		keygenerated bool
	} `json:"Agentinfo"`
	Workstationinfo struct {
		IPaddr       string `json:"WSIP"`
		Port         string `json:"WSport"`
		SSH_username string `json:"SSH username"`
	} `json:"Workstationinfo"`

	Detectionmarker struct {
		Markerisdetected bool
	} `json:"Detection Marker"`
	Filepath struct {
		Configfilepath   string `json:"config file path"`
		Databasefilepath string `json:"database file path"`
		Logfilepath      string `json:"log file path"`
	} `json:"filepath"`
	Backup struct {
		SnapshotNum  int32 `json:"Number of snapshots"`
		FullSnapshot bool  `json:"First full Backup named Ryan checker"`
	} `json:"backup"`
}

func InitializeJSON() error {

	defaultConfig := Config{
		Agentinfo: struct {
			Ipaddr string `json:"AGIP"`
			Port   string `json:"AGport"`
			keygenerated bool
		}{Ipaddr: "", Port: "1969", keygenerated: false},
		Workstationinfo: struct {
			IPaddr       string `json:"WSIP"`
			Port         string `json:"WSport"`
			SSH_username string `json:"SSH username"`
		}{IPaddr: "", Port: "1969", SSH_username: "none"},
		Detectionmarker: struct{ Markerisdetected bool }{Markerisdetected: false},
		Filepath: struct {
			Configfilepath   string `json:"config file path"`
			Databasefilepath string `json:"database file path"`
			Logfilepath      string `json:"log file path"`
		}{Configfilepath: "/etc/ASRS_agent/.config/config.json",
			Databasefilepath: "/etc/ASRS_agent/.database/database.json",
			Logfilepath:      "/etc/ASRS_agent/.database/logs.json"},
		Backup: struct {
			SnapshotNum  int32 `json:"Number of snapshots"`
			FullSnapshot bool  `json:"First full Backup named Ryan checker"`
		}{FullSnapshot: false, SnapshotNum: 0},
	}

	jsonData, err := json.MarshalIndent(defaultConfig, "", "  ")
	errorhandler(err, red+"CONFIG ERROR:  Error parsing config file"+reset)
	err = ioutil.WriteFile(defaultConfig.Filepath.Configfilepath, jsonData, 0766)
	return nil
}

func SSH_config() string {
	filedata, err := ioutil.ReadFile("/etc/ASRS_agent/.config/config.json")
	errorhandler(err, red+"Failed to to read file for SSH config"+reset)
	//ip, _ := WSInfoParser()
	var SSHuser Config
	
	_ = json.Unmarshal(filedata, &SSHuser)
	//user := strings.TrimSpace(string(SSHuser.Workstationinfo.SSH_username))
	//if user != "none" {}
	//userANDip := fmt.Sprintf("%v@%v", user, ip)
	_,err = os.Stat("/etc/ASRS_agent/.config/")
	if os.IsNotExist(err) {
		cmd1 := exec.Command("sudo", "ssh-keygen", "-t", "rsa", "-f", "/etc/ASRS_agent/.config/id_rsa.pub", "-N", `""`)
		// Get a file descriptor for stdin
		cmdout1, err := cmd1.CombinedOutput()
		
		fmt.Println(red+"SSH MESSAGE: Failed to generate SSH keys"+reset,err)
		fmt.Println(string(cmdout1))
	} else if err != nil {
		fmt.Println("SSH MESSAGE: Error checking file:", err)
	}
	keys, err := ioutil.ReadFile("/etc/ASRS_agent/.config/id_rsa.pub")
	errorhandler(err, red+"SSH MESSAGE: keys not found"+reset)
	return string(keys)

	//cmd2 := exec.Command("sudo", "ssh-copy-id", "-i", "/etc/ASRS_agent/.config/id_rsa.pub", userANDip)
	//output_full, err := cmd2.CombinedOutput()
	//errorhandler(err, "Failed to get output of copying keys to workstation")
	//fmt.Println("final : ", string(output_full))
}

// create info parser for whole infos

func configparser() {
	filedata, err := ioutil.ReadFile("/etc/ASRS_agent/.config/config.json")
	if err != nil {
		fmt.Println(red+"Error:"+reset, err)
		return
	}

	var config Config
	err = json.Unmarshal(filedata, &config)
	if err != nil {
		fmt.Println(red+"Error:"+reset, err)
		return
	}

}

func WSInfoParser() (ip, port string) {
	filedata, err := ioutil.ReadFile("/etc/ASRS_agent/.config/config.json")
	if err != nil {
		fmt.Println(red+"Error:"+reset, err)
		return
	}

	var info Config
	err = json.Unmarshal(filedata, &info)
	if err != nil {
		fmt.Println(red+"Error:"+reset, err)
		return
	}
	return info.Workstationinfo.IPaddr, info.Agentinfo.Port

}

func errorhandler(err error, s string) {
	if err != nil {
		fmt.Println(red+"Error: "+reset, s, err)
	}
}
