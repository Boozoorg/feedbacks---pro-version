package models

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
)

type Config struct {
	AppConf          App         `json:"app"`
	SmsConf          SmsConf     `json:"sms"`
	DB               string      `json:"db"`
	TokenAuth        string      `json:"tokenAuth"`
	UrlWtInterchange string      `json:"url_wt_interchange"`
	UrlWtApi         string      `json:"url_wt_api"`
	Logger           Logs        `json:"logs"`
}

type App struct {
	ServerName string `json:"serverName"`
	Port       int64  `json:"portRun"`
	Debug      bool   `json:"debug"`
}

type SmsConf struct {
	Uri string `json:"uri"`
}

type Logs struct {
	Logfile    string `json:"logfile"`
	MaxSize    int    `json:"max_size"`
	MaxBackups int    `json:"max_backups"`
	MaxAge     int    `json:"max_age"`
}

var Conf Config

func Setup(F string) {
	byteValue, err := ioutil.ReadFile(F)
	if err != nil {
		log.Fatalf("%v", err)
		return
	}
	err = json.Unmarshal(byteValue, &Conf)

	fmt.Println(Conf)
	if err != nil {
		log.Fatalf("%v", err)
		return
	}
}
