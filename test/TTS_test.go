package test

import (
	service "../service"
	"fmt"
	"os"
	"testing"
	"time"
)

//This file tests the text to speech functionality of the project.

//TestCredentials ensures that the user has supplied the required credential file.
//Future Implementations of this test would ensure that the files supplied are valid GCP API service account keys,
//but this MVP will expect anyone using the system to understand how to manage API keys.
func TestCredentials(t *testing.T) {
	if os.Getenv("GOOGLE_APPLICATION_CREDENTIALS") == "" {
		fmt.Println("Test Credentials] FAILED - No Google credentials environment export found")
		t.FailNow()
	} else {
		fmt.Println("Test Credentials] PASSED - Google credentials environment export found")
	}
}

//TestPlainTextTTSGeneration is a function which makes a text to speech request to the GCP API and
//specifies that the request be plain text, and not SSML.
func TestPlainTextTTSGeneration(t *testing.T) {

	if !service.CheckForAudioFile("Please Wait For Unit Tests.wav") {
		r := service.SpeechRequest{
			Text:         "Please Wait For Unit Tests",
			LanguageCode: "en-US",
			SsmlGender:   "FEMALE",
			VoiceName:    "en-us-Wavenet-C",
		}

		r.SpeakToFile("../audio/Please Wait For Unit Tests.wav")
		if !service.CheckForAudioFile("Please Wait For Unit Tests.wav") {
			fmt.Println("Attempted to generate audio, OK doing generation but could not find saved audio file ")
			t.FailNow()
		}
	}

	fmt.Println("Plain Text TTS Unit Tests PASSED")
}

func TestSpeakNow(t *testing.T) {
	for i := 0; i < 10; i++ {
		service.SpeakAloud("please wait for unit tests")
		time.Sleep(2 * time.Second)
	}
}
