package main

import (
	"../libgolb"
	"github.com/docopt/docopt-go"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"strings"
)
var origins = map[string]string{}
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

	
	//serv := strings.Split(req.RemoteAddr, ":") // extract just IP without port
	origin := req.RemoteAddr
	libgolb.Log("misc", "Access From :"+origin)
	server, errGS := origins[origin]
	if errGS == false {
		server = getServer()
	}
	limit := 0
	for limit < libgolb.NumberBack {
		resp, _ := http.NewRequest(req.Method, "http://"+server+"/", nil)
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
			defer secondResp.Body.Close()
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
	origins[origin] = server
	libgolb.Log("ok", "Answer From :"+origin)
	libgolb.LogW3C(w, req, false)
}

func parseArgument(configuration string) {

	// Load configuration
	libgolb.ConfLoad(configuration)
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
  golb <configuration>
  golb -h | --help
  golb --version

Options:
  -h --help     Show this screen.
  --version     Show version.`

	arguments, _ := docopt.Parse(usage, nil, true, "GoLB 0.1", false)
	parseArgument(arguments["<configuration>"].(string))
}
