package libgolb

import (
	"io"
	"net/http"
	"strconv"
	"strings"
)

var ()

func HttpResponse(w http.ResponseWriter, status int, message string) string {
	w.Header().Set("Status", strconv.Itoa(status))
	io.Copy(w, strings.NewReader(message))
	return "HTTP Response (" + strconv.Itoa(status) + ") " + message
}
