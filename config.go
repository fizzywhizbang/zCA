package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type ZcaConfig struct {
	RootDIR string `json:"root-dir"`
	CertDir string `json:"cert-dir"`
}

// check if a file exists used for startup
func Exists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

// check for configurations
func CkConfig() bool {
	homedir, err := os.UserHomeDir()
	if err != nil {
		log.Println("Unable to get home dir")
		return false
	}
	configFile := "zca-config.json"
	configDir := homedir + "/.config/zca"
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		//if not exist create it
		err = os.Mkdir(configDir, 0755)
		if err != nil {
			log.Println("Error creating config dir (line 38)")
			return false
		}
		if !Exists(configDir + "/" + configFile) {
			root := configDir + "/root"
			cert := configDir + "/crt"
			return writeConfig(root, cert, configDir, configFile)
		}
		return true
	}
	return true
}

func writeConfig(root, cert, configDir, configFile string) bool {

	file, err := os.Create(configDir + "/" + configFile)
	if err != nil {
		log.Println("Unable to create config file (line 60)")
		return false
	}
	defer file.Close()
	fmt.Fprintln(file, "{")
	fmt.Fprintln(file, "\t\"root-dir\":\""+root+"\",")
	fmt.Fprintln(file, "\t\"cert-dir\":\""+cert+"\",")
	fmt.Fprintln(file, "}")

	return true
}

func ConfigParser() ZcaConfig {
	var config ZcaConfig
	homedir, err := os.UserHomeDir()
	if err != nil {
		log.Println("Unable to get home dir")
	}
	configFile := "zca-config.json"
	configDir := homedir + "/.config/zca"

	configFileFile, err := os.Open(configDir + "/" + configFile)
	if err != nil {
		log.Println("Unable to read config file (line 59)")
	}

	jsonParser := json.NewDecoder(configFileFile)
	err = jsonParser.Decode(&config)
	if err != nil {
		log.Fatal("Can't decode your json", err)
	}
	defer configFileFile.Close()

	return config
}
