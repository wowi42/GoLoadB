package libgolb

import ()

func ErrCatcher(message string, err error) bool {
	if err != nil {
		Log("error", message+" : "+err.Error())
		return false
	}
	return true
}
