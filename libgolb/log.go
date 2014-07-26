package libgolb

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

// ANSIC       = "Mon Jan _2 15:04:05 2006"
// UnixDate    = "Mon Jan _2 15:04:05 MST 2006"
// RubyDate    = "Mon Jan 02 15:04:05 -0700 2006"
// RFC822      = "02 Jan 06 15:04 MST"
// RFC822Z     = "02 Jan 06 15:04 -0700" // RFC822 with numeric zone
// RFC850      = "Monday, 02-Jan-06 15:04:05 MST"
// RFC1123     = "Mon, 02 Jan 2006 15:04:05 MST"
// RFC1123Z    = "Mon, 02 Jan 2006 15:04:05 -0700" // RFC1123 with numeric zone
// RFC3339     = "2006-01-02T15:04:05Z07:00"
// RFC3339Nano = "2006-01-02T15:04:05.999999999Z07:00"
// Kitchen     = "3:04PM"
// // Handy time stamps.
// Stamp      = "Jan _2 15:04:05"
// StampMilli = "Jan _2 15:04:05.000"
// StampMicro = "Jan _2 15:04:05.000000"
// StampNano  = "Jan _2 15:04:05.000000000"

var (
	normalColor = "\033[0m"
	redColor    = "\033[1;31m"
	blueColor   = "\033[1;96m"
	greenColor  = "\033[1;92m"
	yellowColor = "\033[1;33m"
	layout      = time.RFC1123Z
)

func LogFormat(typeL, msg string) (logToWrite, logToDisplay string) {
	// Color
	color := normalColor
	if Conf.LogColor {
		if typeL == "error" {
			color = redColor
		} else if typeL == "misc" {
			color = blueColor
		} else if typeL == "warning" {
			color = yellowColor
		} else if typeL == "ok" {
			color = greenColor
		}
	}

	// Time
	t := fmt.Sprintf("%s", time.Now().Format(layout))

	// Format message
	logToWrite = fmt.Sprintf("[%s] %s\n", t, msg)
	logToDisplay = fmt.Sprintf("%s[%s] %s%s%s", normalColor, t, color, msg, normalColor)

	return
}

func Log(typeL, msg string) {
	// Format log lign
	logW, logD := LogFormat(typeL, msg)

	// Display log
	fmt.Println(logD)

	// Write log
	fileToWrite := ""
	if typeL == "error" {
		fileToWrite = Conf.Log.Folder + Conf.Name + "-error.log"
	} else {
		fileToWrite = Conf.Log.Folder + Conf.Name + ".log"
	}

	// Create logfile
	err := os.MkdirAll(filepath.Dir(fileToWrite), 0700)
	if err != nil {
		log.Println("Can't create log dir ("+filepath.Dir(fileToWrite)+") failed :", err)
		os.Exit(1)
	}
	exec.Command("touch", fileToWrite).Output()

	// Append log to logfile
	f, err := os.OpenFile(fileToWrite, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		log.Println("Can't open log file (" + fileToWrite + ") failed.")
		os.Exit(1)
	}
	defer f.Close()
	if _, err = f.WriteString(logW); err != nil {
		panic(err)
	}
}

func LogW3C(w http.ResponseWriter, req *http.Request, before bool) {
	// Remotehost rfc931 authuser [date] "request" status bytes
	var msg string
	if before == true {
		msg = req.RemoteAddr + " [" + time.Now().Format(layout) + "] \"" + req.Method + "\" " + req.Header.Get("status") + " " + req.Header.Get("Content-Length") + "\n"
	} else {
		msg = req.RemoteAddr + " - [" + time.Now().Format(layout) + "] \"" + req.Method + " " + req.RequestURI + " " + req.Proto + "\" " + w.Header().Get("status") + " " + w.Header().Get("Content-Length") + "\n"
	}
	// Write log
	fileToWrite := Conf.Log.Folder + Conf.Name + ".w3c.log"

	// Create logfile
	err := os.MkdirAll(filepath.Dir(fileToWrite), 0700)
	if err != nil {
		log.Println("Can't create log dir ("+filepath.Dir(fileToWrite)+") failed :", err)
		os.Exit(1)
	}
	exec.Command("touch", fileToWrite).Output()

	// Append log to logfile
	f, err := os.OpenFile(fileToWrite, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		log.Println("Can't open log file (" + fileToWrite + ") failed.")
		os.Exit(1)
	}
	defer f.Close()
	if _, err = f.WriteString(msg); err != nil {
		panic(err)
	}
}
