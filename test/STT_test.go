package test

import (
	"../service"
	"log"
	"testing"
)

func TestRecognize(t *testing.T) ***REMOVED***
	log.Println("Starting recognition !PLEASE SPEAK TO CONTINUE!")
	output, confidence, err := service.Recognize()
	if err != nil ***REMOVED***
		log.Println(err)
		log.Println("Failed to recognize speech.")
		t.FailNow()
	***REMOVED***

	log.Printf("Recognize test PASSED with an output of %s and a confidence of %f", output, confidence)
***REMOVED***
