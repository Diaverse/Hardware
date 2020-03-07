package service

import "os"

//checkForFile
func CheckForFile(fileDir string) bool {
	_, err := os.Stat(fileDir)
	if err != nil {
		return false
	}
	return true
}
