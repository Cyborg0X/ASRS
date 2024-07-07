package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

/*
	type procedures struct {
		pr_get_status struct {
			Proceduree string `json:"procedure"`
			Username   string `json:"SSH username"`
		} `json:"pr_get_status"`
		pr_heal_now struct {
			Proceduree string `json:"procedure"`
			Username   string `json:"SSH username"`
		}

}
*/

var red = "\033[31m"
var reset = "\033[0m"

type Config struct {
	Agentinfo struct {
		Ipaddr             string `json:"AGIP"`
		Port               string `json:"AGport"`
		SSH_Agent_username string `json:"Agent SSH username"`
	} `json:"Agentinfo"`
	Workstationinfo struct {
		IPaddr  string `json:"WSIP"`
		Port    string `json:"WSport"`
		SendSSH bool   `json:"WS SSH username sent"`
	} `json:"Workstationinfo"`

	Detectionmarker struct {
		Markerisdetected bool
	} `json:"Detection Marker"`
	Filepath struct {
		Configfilepath   string `json:"config file path"`
		Databasefilepath string `json:"database file path"`
		Logfilepath      string `json:"log file path"`
	} `json:"filepath"`
}

func InitializeJSON() error {

	defaultConfig := Config{
		Agentinfo: struct {
			Ipaddr             string `json:"AGIP"`
			Port               string `json:"AGport"`
			SSH_Agent_username string `json:"Agent SSH username"`
		}{Ipaddr: "", Port: "1969", SSH_Agent_username: "none"},
		Workstationinfo: struct {
			IPaddr  string `json:"WSIP"`
			Port    string `json:"WSport"`
			SendSSH bool   `json:"WS SSH username sent"`
		}{IPaddr: "", Port: "1969", SendSSH: false},
		Detectionmarker: struct{ Markerisdetected bool }{Markerisdetected: false},
		Filepath: struct {
			Configfilepath   string `json:"config file path"`
			Databasefilepath string `json:"database file path"`
			Logfilepath      string `json:"log file path"`
		}{Configfilepath: "/etc/ASRS_WS/.config/config.json", Databasefilepath: "/etc/ASRS_WS/.database/database.json", Logfilepath: "/etc/ASRS_WS/.database/logs.json"},
	}

	jsonData, err := json.MarshalIndent(defaultConfig, "", "  ")
	errorhandler(err, "Error:  Error parsing config file")
	err = ioutil.WriteFile(defaultConfig.Filepath.Configfilepath, jsonData, 0766)
	errorhandler(err, "Error writing config file:")
	return nil
}

type Logs struct {
	Agentlogs struct {
		SSH_sent bool `json:"marker for username"`
	}
}

func AgentInfoParser() (ip, port string) {
	filedata, err := ioutil.ReadFile("/etc/ASRS_WS/.config/config.json")
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
	return info.Workstationinfo.IPaddr, info.Workstationinfo.Port

}

func errorhandler(err error, s string) {
	if err != nil {
		fmt.Println(red+"Error: "+reset, s, err)
	}
}
