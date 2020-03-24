package controller

import (
	"../domain"
	"../service"
	"encoding/json"
	"fmt"
	"github.com/prometheus/common/log"
	"time"
)

//ExecuteTestScript is the function which executes the core logic of the project, it utilizes the code within service directory of the project to do the required audio I/O.

func ExecuteTestScript(script *domain.TestScript) error {

	for _, e := range script.TestCases {
		testCaseResults := make(map[int]bool)
		for i := 0; i < len(e.HardwareInput); i++ {
			time.Sleep(2 * time.Second) //Any VUI will have some amount of processing time, this value is temporary. We have a technical requirement for this pause as well, as the audio device cannot open and close as fast as this loop
			service.SpeakAloud(e.HardwareOutput[i])

			response, confidence, err := service.Recognize()
			if err != nil {
				return err
			}

			log.Infof("Recognized Response: %s | Confidence of %f", response, confidence)
			if response == e.HardwareInput[i] {
				testCaseResults[i] = true
				e.TotalPassed++
			} else {
				testCaseResults[i] = false
				e.TotalFailed++
			}
		}
		tpass := 0
		tfail := 0
		for _, v := range testCaseResults {
			if v {
				tpass++
			} else {
				tfail++
			}
		}
		fmt.Println(tpass)
		fmt.Println(tfail)
		fmt.Println(testCaseResults)
		if float64(tfail+tpass) != 0 {
			e.Result = float64(tpass) / float64(tfail+tpass)
		} else {
			fmt.Println("Almost divided by zero")
		}
	}

	tpercent := 0.

	for _, e := range script.TestCases {
		tpercent += e.Result
	}

	script.PassPercent = tpercent / float64(len(script.TestCases))
	j, _ := json.Marshal(script)
	fmt.Println(string(j))
	return nil
}
