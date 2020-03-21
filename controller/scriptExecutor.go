package controller

import (
	"../domain"
	"../service"
	"github.com/prometheus/common/log"
	"time"
)

//ExecuteTestScript is the function which executes the core logic of the project, it utilizes the code within service directory of the project to do the required audio I/O.

func ExecuteTestScript(script *domain.TestScript) error {

	for l, e := range script.TestCases {
		for i := 0; i < len(e.HardwareInput); i++ {
			time.Sleep(2 * time.Second) //Any VUI will have some amount of processing time, this value is temporary. We have a technical requirement for this pause as well, as the audio device cannot open and close as fast as this loop
			service.SpeakAloud(e.HardwareOutput[i])

			response, confidence, err := service.Recognize()
			if err != nil {
				return err
			}

			log.Infof("Recognized Response: %s | Confidence of %f", response, confidence)

			if service.TranscriptionConfidence(response, e.HardwareInput[i]) >= .60 {
				log.Infof("Response %d for Test Case %d PASSED", i, l)
				e.TotalPassed++
			} else {
				log.Infof("Response %d for Test Case %d FAILED", i, l)
				e.TotalFailed++
			}
		}
		e.Result = float64(e.TotalPassed/e.TotalFailed + e.TotalPassed)
	}

	TotalResult := 0.
	//process results.
	for _, e := range script.TestCases {
		TotalResult += e.Result
	}

	TotalResult = TotalResult / float64(len(script.TestCases))

	return nil
}
