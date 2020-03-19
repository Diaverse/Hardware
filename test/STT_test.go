package test

import (
	"../service"
	"log"
	"testing"
)

func TestRecognize(t *testing.T) ***REMOVED***
	output, confidence, err := service.Recognize()
	if err != nil ***REMOVED***
		log.Println("Failed to recognize speech.")
		t.FailNow()
	***REMOVED***

	log.Printf("Recognize test passed with an output of %s and a confidence of %f", output, confidence)
***REMOVED***
