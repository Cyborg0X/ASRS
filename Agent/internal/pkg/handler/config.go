package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type Config struct {
	Agentinfo struct {
		Ipaddr string `json:"AGIP"`
		Port   string `json:"AGport"`
	} `json:"Agentinfo"`
	Workstationinfo struct {
		IPaddr string `json:"WSIP"`
		Port   string `json:"WSport"`
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
			Ipaddr string `json:"AGIP"`
			Port   string `json:"AGport"`
		}{Ipaddr: "", Port: "1969"},
		Workstationinfo: struct {
			IPaddr string `json:"WSIP"`
			Port   string `json:"WSport"`
		}{IPaddr: "", Port: "1969"},
		Detectionmarker: struct{ Markerisdetected bool }{},
		Filepath: struct {
			Configfilepath   string `json:"config file path"`
			Databasefilepath string `json:"database file path"`
			Logfilepath      string `json:"log file path"`
		}{Configfilepath: "/etc/ASRS_agent/.config/config.json", Databasefilepath: "/etc/ASRS_agent/.database/database.json", Logfilepath: "/etc/ASRS_agent/.database/logs.json"},
	}

	jsonData, err := json.Marshal(defaultConfig)
	errorhandler(err, "Error:  Error parsing config file")
	err = ioutil.WriteFile(defaultConfig.Filepath.Configfilepath, jsonData, 0644)
	return nil
}

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



func errorhandler(err error, s string) {
	if err != nil {
		fmt.Println("Error: ", s, err)
	}
}
