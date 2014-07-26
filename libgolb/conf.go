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

type RedisServer struct {
	Hostname string
	Port     string
	Database int
}

type HttpServer struct {
	Hostname string
	Port     string
}

type Configuration struct {
	Name         string
	Log          LogFile
	RedisLB      RedisServer
	Server       HttpServer
	BackServers  []string
	LogColor     bool
	TTL          int
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
}
