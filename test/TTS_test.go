package test

import (
	service "../service"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"testing"
)

//This file tests the text to speech functionality of the project.

//TestCredentials ensures that the user has supplied the required credential file.
//Future Implementations of this test would ensure that the files supplied are valid GCP API service account keys,
//but this MVP will expect anyone using the system to understand how to manage API keys.
func TestCredentials(t *testing.T) ***REMOVED***
	files, err := ioutil.ReadDir("../credentials")
	if err != nil ***REMOVED***
		log.Fatal(err)
	***REMOVED***

	jsonFileFound := false

	for _, e := range files ***REMOVED***
		if strings.Contains(e.Name(), ".json") ***REMOVED***
			jsonFileFound = true
		***REMOVED***
	***REMOVED***

	if !jsonFileFound ***REMOVED***
		fmt.Println("ERROR - could not find credentials file with proper naming conventions. Ensure your GCP API keys are placed  within the credentials directory and  end in '.json'.")
		t.FailNow()
	***REMOVED***
	log.Println("Credentials Are Valid PASSED")
***REMOVED***

//TestPlainTextTTSGeneration is a function which makes a text to speech request to the GCP API and
//specifies that the request be plain text, and not SSML.
func TestPlainTextTTSGeneration(t *testing.T) ***REMOVED***

	r := service.SpeechRequest***REMOVED***
		Text:         "Please Wait For Unit Tests",
		LanguageCode: "en-US",
		SsmlGender:   "FEMALE",
		VoiceName:    "en-us-Wavenet-C",
	***REMOVED***

	r.SpeakToFile("../audio/Please Wait For Unit Tests.wav")
	if !service.CheckForAudioFile("Please Wait For Unit Tests.wav") ***REMOVED***
		fmt.Println("Could not find audio file.")
		t.FailNow()
	***REMOVED***
	fmt.Println("Plain Text TTS Unit Tests PASSED")
***REMOVED***

//TestSsmlTTSGeneration is a function which makes a text to speech request to the GCP API and
//specifies that the request be plain text, and not SSML.
//func TestSsmlTTSGeneration(t *testing.T) ***REMOVED***
//	service.SpeakAloud("<speak><prosody  pitch=\"-20%\" range=\"low\" rate =\"50%\" volume =\"30\"><emphasis level=\"strong\">Almost done</emphasis></prosody></speak>")
//	_, e := os.Stat("Almost done.wav")
//	if e != nil ***REMOVED***
//		log.Println("Could not generate audio file for TTS unit test")
//		t.FailNow()
//	***REMOVED***
//***REMOVED***
