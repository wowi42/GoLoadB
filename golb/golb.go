package main

import (
	"../libgolb"
	"github.com/docopt/docopt-go"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"strings"
	"net"
	"time"
	"os"
)

var origins = map[string]string{} // map of originip/backend
var redisArg bool

func GetAndCheckServer() (server string) {
	limit := 0
	for limit < libgolb.NumberBack {
		server = libgolb.Conf.BackServers[libgolb.RoundRobin]
		connection, err := net.DialTimeout("tcp", server, time.Duration(libgolb.Conf.TimeOut) * time.Second)
		if err == nil {
			connection.Close()
			return
		}
		libgolb.RoundRobin++
		if libgolb.RoundRobin >= libgolb.NumberBack {
			libgolb.RoundRobin = 0
		}
	}
	return
}

func extractKey(RemoteAddr string) (key string){
	if libgolb.Conf.IpHashLevel != 5 {
		ip := strings.Split(RemoteAddr, ":")
		octets := strings.Split(ip[0], ".")		
		key = strings.Join(octets[:libgolb.Conf.IpHashLevel], ".")
	} else {
		key = RemoteAddr
	}
	return
	
}

func GetAddress(origin string) (server string, result bool){
	var err error
	if redisArg == true {
		server, err = libgolb.RadixGetString(libgolb.LBClient, origin)
		if err != nil {
			result = false
		} else {
			result = true
		}
	} else {
		server, result = origins[origin]
	}
	return
}

func golbGet(w http.ResponseWriter, req *http.Request) {
	var secondResp *http.Response
	var errsp error

	origin := extractKey(req.RemoteAddr)
	libgolb.Log("misc", "Access From: "+ req.RemoteAddr + " Key: " +origin)
	server, errGS := GetAddress(origin) 
	if errGS == false {
		server = GetAndCheckServer()
	}
	primaryServer := server
	limit := 0
	for limit < libgolb.NumberBack { // this for is used to check all servers and select the first one available
		resp, _ := http.NewRequest(req.Method, "http://"+server+"/", nil)
		resp.Header = req.Header
		resp.Header.Set("X-Forwarded-For", req.RemoteAddr)
		secondResp, errsp = http.DefaultClient.Do(resp)
		if errsp != nil {
			libgolb.Log("error", "Connection with the HTTP file server failed: "+errsp.Error())
			server = GetAndCheckServer()
			limit++
		} else {
			defer secondResp.Body.Close() // don't forget to close the Body !!!
			break
		}
	}
	if limit >= libgolb.NumberBack { // No Backend
		libgolb.HttpResponse(w, 500, "Internal server error\n")
		libgolb.Log("error", "No Backend Server avalaible")
		return
	}
	for k, v := range secondResp.Header { // Copy Header
		w.Header().Add(k, strings.Join(v, ""))
	}
	w.Header().Add("Served-By", server)
	w.Header().Add("Server", libgolb.Conf.Name)	
	io.Copy(w, secondResp.Body)
	if redisArg == false {
		if primaryServer != server {
			origins[origin] = server
		}
	} else {
		if primaryServer != server {
			_ = libgolb.RadixSet(libgolb.LBClient, origin, server)
		}
		_ = libgolb.RadixExpire(libgolb.LBClient, origin)
	}
	libgolb.Log("ok", "Answer From :"+server)
	libgolb.LogW3C(w, req, false)
}

func parseArgument(configuration string) {
	// Load configuration
	libgolb.ConfLoad(configuration)
	//Connect to Redis
	redis := libgolb.ConnectToRedis()
	if redis != nil {
		libgolb.Log("error", "Redis connection failed: Server = "+libgolb.Conf.RedisLB.Hostname+":"+libgolb.Conf.RedisLB.Port)
		os.Exit(1)
	}

	// Router
	rtr := mux.NewRouter()
	rtr.HandleFunc("/", golbGet).Methods("GET")
	http.Handle("/", rtr)

	// Listening
	libgolb.Log("ok", "Listening on "+libgolb.Conf.Server.Hostname+":"+libgolb.Conf.Server.Port)
	err := http.ListenAndServe(libgolb.Conf.Server.Hostname+":"+libgolb.Conf.Server.Port, nil)
	libgolb.ErrCatcher("ListenAndServe: ", err)
}

func main() {
	usage := `Golb.

Usage:
  golb memory <configuration>
  golb redis <configuration>
  golb -h | --help
  golb --version

Options:
  -h --help     Show this screen.
  --version     Show version.`

	arguments, _ := docopt.Parse(usage, nil, true, "GoLB 0.1", false)
	if arguments["redis"] == true {
		redisArg = true
	} else {
		redisArg = false
	}
	parseArgument(arguments["<configuration>"].(string))
}
