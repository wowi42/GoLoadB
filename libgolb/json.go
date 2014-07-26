package libgolb

import (
	"encoding/json"
)

func ServerEncode(f Server) string {
	s, _ := json.Marshal(f)
	return string(s[:])
}

func ServerDecode(s string) Server {
	var f Server
	_ = json.Unmarshal([]byte(s), &f)
	return f
}
