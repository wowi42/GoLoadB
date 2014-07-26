package main

import (
	"../libgolb"
	"github.com/docopt/docopt-go"
	"github.com/gorilla/mux"
	"net/http"
	"os"
	"strings"
	"io"
	"strconv"
)
func getServer() (server string) {
	server = libgolb.Conf.BackServers[libgolb.RoundRobin]
	libgolb.RoundRobin++
	if libgolb.RoundRobin >= libgolb.NumberBack {
		libgolb.RoundRobin = 0
	}
	return
}

func golbGet(w http.ResponseWriter, req *http.Request) {
	var secondResp *http.Response
	var errsp error
	
	serv := strings.Split(req.RemoteAddr, ":") // extract just IP without port
	libgolb.Log("misc", "Access From :"+serv[0])
	server, errGS := libgolb.RadixGetString(libgolb.LBClient, serv[0])
	if errGS != nil {
		server = getServer()
	}
	limit := 0
	for limit < libgolb.NumberBack {
		resp, _ := http.NewRequest(req.Method, "http://"+server, nil)
		for k, v := range req.Header {
			resp.Header[k] = v
		}
		resp.Header.Set("X-Forwarded-For", req.RemoteAddr)
		secondResp, errsp = http.DefaultClient.Do(resp)
		if errsp != nil {
			libgolb.Log("error", "Connection with the HTTP file server failed: "+errsp.Error())
			server = getServer()
			limit++
		} else {
			break
		}
	}
	if limit >= libgolb.NumberBack {
		libgolb.HttpResponse(w, 500, "Internal server error\n")
		libgolb.Log("error", "No Backend Server avalaible")
		return
	}
	for k, v := range secondResp.Header {
		w.Header().Add(k, strings.Join(v, ""))
	}
	w.Header().Set("Status", "200")
	io.Copy(w, secondResp.Body)
	libgolb.RadixSet(libgolb.LBClient, serv[0], server)
	libgolb.Log("ok", "Answer From :"+serv[0])
	//TTL
	libgolb.RadixExpire(libgolb.LBClient, serv[0], strconv.Itoa(libgolb.Conf.TTL))
	libgolb.LogW3C(w, req, false)
}

func parseArgument(configuration string) {

	// Load configuration
	libgolb.ConfLoad(configuration)
	// Check Redis connection
	redis := libgolb.ConnectToRedis()
	if redis != nil {
		libgolb.Log("error", "Redis connection failed")
		os.Exit(1)
	}

	// Router
	rtr := mux.NewRouter()
	rtr.HandleFunc(`/{URI}`, golbGet).Methods("GET")
	http.Handle("/", rtr)

	// Listening
	libgolb.Log("ok", "Listening on "+libgolb.Conf.Server.Hostname+":"+libgolb.Conf.Server.Port)
	err := http.ListenAndServe(libgolb.Conf.Server.Hostname+":"+libgolb.Conf.Server.Port, nil)
	libgolb.ErrCatcher("ListenAndServe: ", err)
}

func main() {
	usage := `Golb.

Usage:
  golb <configuration>
  golb -h | --help
  golb --version

Options:
  -h --help     Show this screen.
  --version     Show version.`

	arguments, _ := docopt.Parse(usage, nil, true, "GoLB 0.1", false)
	parseArgument(arguments["<configuration>"].(string))
}
