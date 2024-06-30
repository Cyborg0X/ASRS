package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os/exec"
	"strings"
)

type Config struct {
	Agentinfo struct {
		Ipaddr string `json:"AGIP"`
		Port   string `json:"AGport"`
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
		}{Ipaddr: "", Port: "1969"},
		Workstationinfo: struct {
			IPaddr       string `json:"WSIP"`
			Port         string `json:"WSport"`
			SSH_username string `json:"SSH username"`
		}{IPaddr: "", Port: "1969", SSH_username: ""},
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
	errorhandler(err, "CONFIG ERROR:  Error parsing config file")
	err = ioutil.WriteFile(defaultConfig.Filepath.Configfilepath, jsonData, 0766)
	return nil
}

func SSH_config() {
	filedata, err := ioutil.ReadFile("/etc/ASRS_agent/.config/config.json")
	ip, _ := WSInfoParser()
	cmd1 := exec.Command("sudo", "ssh-keygen", "-t", "rsa", "-f", "/etc/ASRS_agent/.config/id_rsa.pub", "-N", `""`)

	// Get a file descriptor for stdin
	cmdout1, err := cmd1.CombinedOutput()
	errorhandler(err, "Failed to generate SSH keys")
	fmt.Println(string(cmdout1))

	var SSHuser Config
	//errorhandler(err, "Failed to get output of whoami command")
	_ = json.Unmarshal(filedata, &SSHuser)
	user := strings.TrimSpace(string(SSHuser.Workstationinfo.SSH_username))
	userANDip := fmt.Sprintf("%v@%v", user, ip)
	cmd2 := exec.Command("sudo", "ssh-copy-id", "-i", "/etc/ASRS_agent/.config/id_rsa.pub", userANDip)
	output_full, err := cmd2.CombinedOutput()
	errorhandler(err, "Failed to get output of copying keys to workstation")
	fmt.Println("final : ", string(output_full))
}

// create info parser for whole infos

func configparser() {
	filedata, err := ioutil.ReadFile("/etc/ASRS_agent/.config/config.json")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	var config Config
	err = json.Unmarshal(filedata, &config)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

}

func WSInfoParser() (ip, port string) {
	filedata, err := ioutil.ReadFile("/etc/ASRS_agent/.config/config.json")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	var info Config
	err = json.Unmarshal(filedata, &info)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	return info.Workstationinfo.IPaddr, info.Agentinfo.Port

}

func errorhandler(err error, s string) {
	if err != nil {
		fmt.Println("Error: ", s, err)
	}
}
