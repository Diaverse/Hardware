package service

import "os"

//checkForFile
func CheckForFile(fileDir string) bool ***REMOVED***
	_, err := os.Stat(fileDir)
	if err != nil ***REMOVED***
		return false
	***REMOVED***
	return true
***REMOVED***
