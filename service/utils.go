package service

import (
	"fmt"
	"os"
)

//checkForFile
func CheckForFile(fileDir string) bool {
	_, err := os.Stat(fileDir)
	if err != nil {
		return false
	}
	return true
}

func CheckForAudioFile(fileName string) bool {
	_, err := os.Stat("/home/pi/Hardware/audio/" + fileName)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}

//credit for below function https://bit.ly/3ae0afW
func DumpMap(space string, m map[string]interface{}) {
	for k, v := range m {
		if mv, ok := v.(map[string]interface{}); ok {
			fmt.Printf("{ \"%v\": \n", k)
			DumpMap(space+"\t", mv)
			fmt.Printf("}\n")
		} else {
			fmt.Printf("%v %v : %v\n", space, k, v)
		}
	}
}
