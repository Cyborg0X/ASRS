package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"
)

type Config struct {
	Agentinfo struct {
		Ipaddr string `json:"AGIP"`
		Port   string `json:"AGport"`
		//keygenerated bool
	} `json:"Agentinfo"`
	Workstationinfo struct {
		IPaddr        string `json:"WSIP"`
		Port          string `json:"WSport"`
		Webuser       string `json:"web files rsync user"`
		SnapshotsUser string `json:"snapshots rsync user"`
		//SSH_username string `json:"SSH username"`
		//SSHpass string `json:"password"`
	} `json:"Workstationinfo"`

	Detectionmarker struct {
		Markerisdetected bool
		AttackerIP       string `json:"attacker IP"`
		AttackTiming     string `json:"Time of attack"`
	} `json:"Detection Marker"`
	Filepath struct {
		Configfilepath   string `json:"config file path"`
		Databasefilepath string `json:"database file path"`
		Logfilepath      string `json:"log file path"`
	} `json:"filepath"`
	Backup struct {
		SnapshotNum  int32  `json:"Number of snapshots"`
		FullSnapshot bool   `json:"First full Backup named Ryan checker"`
		Ltimestamp   string `json:"timestamp of the last backup"`
	} `json:"backup"`
}

func InitializeJSON() error {

	defaultConfig := Config{
		Agentinfo: struct {
			Ipaddr string `json:"AGIP"`
			Port   string `json:"AGport"`
			//keygenerated bool
		}{Ipaddr: "", Port: "1969"},
		Workstationinfo: struct {
			IPaddr        string `json:"WSIP"`
			Port          string `json:"WSport"`
			Webuser       string `json:"web files rsync user"`
			SnapshotsUser string `json:"snapshots rsync user"`
			//SSH_username string `json:"SSH username"`
			//SSHpass string `json:"password"`
		}{IPaddr: "", Port: "1969", Webuser: "webuser", SnapshotsUser: "asrs"},
		Detectionmarker: struct {
			Markerisdetected bool
			AttackerIP       string `json:"attacker IP"`
			AttackTiming     string `json:"Time of attack"`
		}{Markerisdetected: false, AttackerIP: "", AttackTiming: ""},
		Filepath: struct {
			Configfilepath   string `json:"config file path"`
			Databasefilepath string `json:"database file path"`
			Logfilepath      string `json:"log file path"`
		}{Configfilepath: "/etc/ASRS_agent/.config/config.json",
			Databasefilepath: "/etc/ASRS_agent/.database/database.json",
			Logfilepath:      "/etc/ASRS_agent/.database/logs.json"},
		Backup: struct {
			SnapshotNum  int32  `json:"Number of snapshots"`
			FullSnapshot bool   `json:"First full Backup named Ryan checker"`
			Ltimestamp   string `json:"timestamp of the last backup"`
		}{FullSnapshot: false, SnapshotNum: 0, Ltimestamp: ""},
	}
	//fileper, _ := os.Stat(filepath)
	//per := fileper.Mode().Perm()
	jsonData, err := json.MarshalIndent(defaultConfig, "", "  ")
	if err != nil {
		fmt.Println(err, red+"CONFIG ERROR:  Error parsing config file:"+reset, err)
	}
	err = ioutil.WriteFile("/etc/ASRS_agent/.config/config.json", jsonData, 0777)
	if err != nil {
		fmt.Println("CONFIG MESSAGE: Failed to write data to config file", err)
	}
	return nil
}

/*
func SSH_config() []byte {
	fmt.Println("SSH CONFIG STARTED")

	info,err := os.Stat("/etc/ASRS_agent/.config/id_rsa.pub")
	if os.IsNotExist(err) {
		cmd1 := exec.Command("sudo", "ssh-keygen", "-t", "rsa", "-f", "/etc/ASRS_agent/.config/id_rsa.pub", "-N", `""`)
		// Get a file descriptor for stdin
		cmdout1, err := cmd1.CombinedOutput()

		fmt.Println(red+"SSH MESSAGE: Failed to generate SSH keys"+reset,err)
		fmt.Println(string(cmdout1))
	} else if err != nil {
		fmt.Println("SSH MESSAGE: Error checking file:", err)
	}
	keys := make([]byte, info.Size())
	file, err := os.Open("/etc/ASRS_agent/.config/id_rsa.pub")
	errorhandler(err, "SSH MESSAGE: Can't open SSH key file")
	defer file.Close()
	_,err = file.Read(keys)
	errorhandler(err, red+"SSH MESSAGE: can't read keys file"+reset)
	return keys
}
*/
// create info parser for whole infos
/*
func configparser() *Config {
	filedata, err := ioutil.ReadFile("/etc/ASRS_agent/.config/config.json")
	if err != nil {
		fmt.Println(red+"Error:"+reset, err)
	}

	var config Config
	err = json.Unmarshal(filedata, &config)
	if err != nil {
		fmt.Println(red+"Error:"+reset, err)
	}
	return &config

}
*/

func WSInfoParser(er chan string) (ip, port string) {
	filedata, err := ioutil.ReadFile("/etc/ASRS_agent/.config/config.json")
	Errorhandler(err, "WS INFOPARSER MESSAGE: failed to read data", er)
	var info Config
	err = json.Unmarshal(filedata, &info)
	Errorhandler(err, "WS INFOPARSER MESSAGE: failed to Unmarshal data", er)
	return info.Workstationinfo.IPaddr, info.Agentinfo.Port

}

func Errorhandler(err error, s string, erro chan string) {
	if err != nil {
		g := fmt.Sprintf("%v: %v", s, err)
		//ioutil.WriteFile("/etc/ASRS_agent/.config/error.txt",[]byte(g), 0755)
		time.Sleep(time.Second *1)
		erro <- g
	}

}
func EventHandler(s string, eve chan string) {
	//ioutil.WriteFile("/etc/ASRS_agent/.config/event.txt",[]byte(s), 0755)
	time.Sleep(time.Second *1)
	eve <- s
}

func NotiHandler(s string, noti chan string) {
	//ioutil.WriteFile("/etc/ASRS_agent/.config/noti.txt",[]byte(s), 0755)
	time.Sleep(time.Second *1)
	noti <- s
}

func ProgHandler(s string, prog chan string) {
	//ioutil.WriteFile("/etc/ASRS_agent/.config/progress.txt",[]byte(s), 0755)
	time.Sleep(time.Second *1)
	prog <- s
}
