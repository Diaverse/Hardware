package controller

import (
	"../domain"
	"../service"
	"encoding/json"
	"fmt"
	"github.com/prometheus/common/log"
)

//ExecuteTestScript is the function which executes the core logic of the project, it utilizes the code within service directory of the project to do the required audio I/O.
type TestCase struct {
	HardwareOutput []string `json:"hardwareOutput"`
	HardwareInput  []string `json:"hardwareInput"`
	Result         float64  `json:"-"`
	TotalPassed    int      `json:"totalPass, omitempty"`
	TotalFailed    int      `json:"totalFail, omitempty"`
}

type TestScript struct {
	TestCases   []TestCase `json:"testCases"`
	PassPercent float64    `json:"passPercent, omitempty"`
}

func ExecuteTestScript(script *domain.TestScript) error {

	for _, e := range script.TestCases {
		testCaseResults := make(map[int]bool)
		service.PrepareAudioFiles(e.HardwareOutput)
		for i := 0; i < len(e.HardwareInput); i++ {
			service.SpeakAloud(e.HardwareOutput[i])
			response, confidence, err := service.Recognize()
			if err != nil {
				return err
			}

			log.Infof("Recognized Response: %s | Confidence of %f", response, confidence)
			transcriptionConf := service.TranscriptionConfidence(response, e.HardwareInput[i])
			log.Infof("---> %f <-----", transcriptionConf)
			if transcriptionConf > .7 {
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
		tpassPercent := float64(tpass) / float64(tpass+tfail)
		e.TotalPassed = tpass
		e.TotalFailed = tfail
		e.Result = tpassPercent
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
