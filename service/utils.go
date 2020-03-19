package service

import (
	"fmt"
	"os"
)

//checkForFile
func CheckForFile(fileDir string) bool ***REMOVED***
	_, err := os.Stat(fileDir)
	if err != nil ***REMOVED***
		return false
	***REMOVED***
	return true
***REMOVED***

//credit for below function https://bit.ly/3ae0afW
func DumpMap(space string, m map[string]interface***REMOVED******REMOVED***) ***REMOVED***
	for k, v := range m ***REMOVED***
		if mv, ok := v.(map[string]interface***REMOVED******REMOVED***); ok ***REMOVED***
			fmt.Printf("***REMOVED*** \"%v\": \n", k)
			DumpMap(space+"\t", mv)
			fmt.Printf("***REMOVED***\n")
		***REMOVED*** else ***REMOVED***
			fmt.Printf("%v %v : %v\n", space, k, v)
		***REMOVED***
	***REMOVED***
***REMOVED***
