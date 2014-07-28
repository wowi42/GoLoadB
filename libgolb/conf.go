package libgolb

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

type LogFile struct {
	Folder string
}

type HttpServer struct {
	Hostname string
	Port     string
}

type RedisServer struct {
	Hostname string
	Port     string
	Database int
}

type Configuration struct {
	Name        string
	Log         LogFile
	Server      HttpServer
	BackServers []string
	LogColor    bool
	TTL         int
	TimeOut     int
	IpHashLevel int
	RedisLB     RedisServer
}

var RoundRobin int
var NumberBack int
var Conf Configuration

func ConfLoad(pathconf string) {
	contents, err := ioutil.ReadFile(pathconf)
	if err != nil {
		log.Fatal("The configuration file (" + pathconf + ") doesn't exist.")
		os.Exit(1)
	}
	err = json.Unmarshal([]byte(contents), &Conf)
	if err != nil {
		log.Fatal("The syntax of the configuration file (" + pathconf + ") is incorrect.")
		os.Exit(1)
	}
	RoundRobin = 0
	NumberBack = len(Conf.BackServers)
	if NumberBack == 0 {
		log.Fatal("No Backend Server !!!")
		os.Exit(1)
	}
	if Conf.IpHashLevel < 0 || Conf.IpHashLevel > 5 {
		log.Fatal("Bad Ip Hash Level !!!")
		os.Exit(1)	
	}
}
