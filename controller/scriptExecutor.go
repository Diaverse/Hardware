package controller

import (
	"../domain"
	"../service"
	"github.com/prometheus/common/log"
	"time"
)

//ExecuteTestScript is the function which executes the core logic of the project, it utilizes the code within service directory of the project to do the required audio I/O.

func ExecuteTestScript(script *domain.TestScript) error ***REMOVED***

	for l, e := range script.TestCases ***REMOVED***
		for i := 0; i < len(e.HardwareInput); i++ ***REMOVED***
			time.Sleep(2 * time.Second) //Any VUI will have some amount of processing time, this value is temporary. We have a technical requirement for this pause as well, as the audio device cannot open and close as fast as this loop
			service.SpeakAloud(e.HardwareOutput[i])

			response, confidence, err := service.Recognize()
			if err != nil ***REMOVED***
				return err
			***REMOVED***

			log.Infof("Recognized Response: %s | Confidence of %f", response, confidence)

			if service.TranscriptionConfidence(response, e.HardwareInput[i]) >= .60 ***REMOVED***
				log.Infof("Response %d for Test Case %d PASSED", i, l)
			***REMOVED*** else ***REMOVED***
				log.Infof("Response %d for Test Case %d FAILED", i, l)
			***REMOVED***
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***
